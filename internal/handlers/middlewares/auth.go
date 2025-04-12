package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
)

type Key string

var Secret = []byte(os.Getenv("JWT_SECRET")) // CHECK

var errSigningMethod = errors.New("signing method error")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("auth-x"), "Bearer ")
		if len(authHeader) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			errorDto := &dto.ErrorDto{
				Message: "Токен сформирован неверно",
			}
			err := json.NewEncoder(w).Encode(errorDto)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		jwtToken := authHeader[1]
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errSigningMethod
			}
			return Secret, nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			errorDto := &dto.ErrorDto{
				Message: "Неверный токен",
			}
			err = json.NewEncoder(w).Encode(errorDto)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), Key("props"), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			errorDto := &dto.ErrorDto{
				Message: "Неверный токен",
			}
			err = json.NewEncoder(w).Encode(errorDto)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	})
}
