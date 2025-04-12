package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hamillka/avitoTechSpring25/internal/db"
	"github.com/hamillka/avitoTechSpring25/internal/handlers"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
	"github.com/hamillka/avitoTechSpring25/internal/logger"
	"github.com/hamillka/avitoTechSpring25/internal/repositories"
	"github.com/hamillka/avitoTechSpring25/internal/usecases"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func setupTestEnvironment(t *testing.T) (http.Handler, func()) {
	dbConfig := &db.DatabaseConfig{
		DBHost: "localhost",
		DBPort: "5432",
		DBName: "pvz_service_test",
		DBUser: "postgres",
		DBPass: "postgres",
	}

	testDB, err := setupTestDatabase(dbConfig)
	require.NoError(t, err)

	logConfig := logger.LogConfig{
		Level: zapcore.DebugLevel,
	}
	testLogger := logger.CreateLogger(logConfig)

	pr := repositories.NewProductRepository(testDB)
	pvzr := repositories.NewPVZRepository(testDB)
	rr := repositories.NewReceptionRepository(testDB)
	ur := repositories.NewUserRepository(testDB)

	ps := usecases.NewProductService(pr, rr, pvzr)
	pvzs := usecases.NewPVZService(pvzr, rr, pr)
	rs := usecases.NewReceptionService(pvzr, rr)
	us := usecases.NewUserService(ur)

	router := handlers.Router(ps, pvzs, rs, us, testLogger)

	cleanup := func() {
		err := testDB.Close()
		if err != nil {
			testLogger.Errorf("Error closing test database: %v", err)
		}

		err = testLogger.Sync()
		if err != nil {
			fmt.Printf("Error syncing logger: %v\n", err)
		}
	}

	return router, cleanup
}

func setupTestDatabase(config *db.DatabaseConfig) (*sqlx.DB, error) {
	testDB, err := db.CreateConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return testDB, nil
}

func TestReceptionFlow(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	pvzID := createPVZ(t, router, "Москва")

	receptionID := createReception(t, router, pvzID)

	products := []string{"электроника", "одежда", "обувь"}
	for i := 1; i <= 50; i++ {
		addProduct(t, router, products[i%3], pvzID)
	}
	time.Sleep(3 * time.Second)

	closeReception(t, router, pvzID)

	verifyReceptionClosed(t, router, pvzID, receptionID)
}

func createPVZ(t *testing.T, router http.Handler, city string) string {
	reqBody := dto.CreatePVZRequestDto{
		City: city,
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	token := getAuthToken(t, router, dto.RoleModerator)
	req.Header.Set("auth-x", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusCreated, resp.Code)

	var response dto.CreatePVZResponseDto
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	require.NoError(t, err)

	return response.Id
}

func createReception(t *testing.T, router http.Handler, pvzID string) string {
	reqBody := dto.CreateReceptionRequestDto{
		PVZId: pvzID,
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	token := getAuthToken(t, router, dto.RoleEmployee)
	req.Header.Set("auth-x", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusCreated, resp.Code)

	var response dto.CreateReceptionResponseDto
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	require.NoError(t, err)

	return response.Id
}

func addProduct(t *testing.T, router http.Handler, prodType, pvzId string) {
	reqBody := dto.AddProductRequestDto{
		Type:  prodType,
		PVZId: pvzId,
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	token := getAuthToken(t, router, dto.RoleEmployee)
	req.Header.Set("auth-x", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusCreated, resp.Code)
}

func closeReception(t *testing.T, router http.Handler, pvzID string) {
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/pvz/%s/close_last_reception", pvzID), nil)
	req.Header.Set("Content-Type", "application/json")

	token := getAuthToken(t, router, dto.RoleEmployee)
	req.Header.Set("auth-x", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
}

func verifyReceptionClosed(t *testing.T, router http.Handler, pvzID, receptionID string) {
	req := httptest.NewRequest(http.MethodGet, "/pvz", nil)

	token := getAuthToken(t, router, dto.RoleEmployee)
	req.Header.Set("auth-x", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var pvzData []dto.PVZWithReceptionsDto
	err := json.Unmarshal(resp.Body.Bytes(), &pvzData)
	require.NoError(t, err)

	var foundPVZ bool
	var foundReception bool

	for _, pvz := range pvzData {
		if pvz.PVZ.Id == pvzID {
			foundPVZ = true

			for _, reception := range pvz.Receptions {
				if reception.Reception.Id == receptionID {
					foundReception = true
					assert.Equal(t, "close", reception.Reception.Status, "Reception should be closed")
					assert.Equal(t, 50, len(reception.Products), "Reception should have 50 products")
					break
				}
			}
			break
		}
	}

	assert.True(t, foundPVZ, "The created PVZ should be found in the response")
	assert.True(t, foundReception, "The created reception should be found in the PVZ data")
}

func getAuthToken(t *testing.T, router http.Handler, role string) string {
	reqBody := dto.DummyLoginRequestDto{
		Role: role,
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var response dto.UserLoginResponseDto
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	require.NoError(t, err)

	return response.Token
}
