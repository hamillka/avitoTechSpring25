package repositories

import (
	"database/sql"
	"time"

	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
	"github.com/hamillka/avitoTechSpring25/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ReceptionRepository struct {
	db *sqlx.DB
}

const (
	getLastReception              = "SELECT id, date_time, pvz_id, status FROM receptions WHERE pvz_id = $1 ORDER BY date_time DESC LIMIT 1"
	createReception               = "INSERT INTO receptions (pvz_id) VALUES ($1) RETURNING id, date_time, pvz_id, status"
	updateReceptionStatus         = "UPDATE receptions SET status = $1 WHERE id = $2 RETURNING id, date_time, pvz_id, status"
	getReceptionsByPVZIdsFiltered = `
	SELECT r.id, r.date_time, r.pvz_id, r.status 
	FROM receptions r
	WHERE r.pvz_id = ANY($1) 
	AND EXISTS (
		SELECT 1 FROM products p 
		WHERE p.reception_id = r.id 
		AND p.date_time BETWEEN $2 AND $3
	)
`
	getReceptionsByPVZIds = `
	SELECT id, date_time, pvz_id, status 
	FROM receptions 
	WHERE pvz_id = ANY($1)
`
)

func NewReceptionRepository(db *sqlx.DB) *ReceptionRepository {
	return &ReceptionRepository{
		db: db,
	}
}

func (rr *ReceptionRepository) GetLastReception(pvzId string) (models.Reception, error) {
	var reception models.Reception
	err := rr.db.QueryRow(getLastReception, pvzId).
		Scan(
			&reception.Id,
			&reception.DateTime,
			&reception.PVZId,
			&reception.Status,
		)
	if err != nil {
		return models.Reception{}, dto.ErrDBRead
	}

	return reception, nil
}

func (rr *ReceptionRepository) CreateReception(pvzId string) (models.Reception, error) {
	var reception models.Reception

	err := rr.db.QueryRow(createReception, pvzId).
		Scan(
			&reception.Id,
			&reception.DateTime,
			&reception.PVZId,
			&reception.Status,
		)
	if err != nil {
		return models.Reception{}, dto.ErrDBInsert
	}

	return reception, nil
}

func (rr *ReceptionRepository) UpdateReceptionStatus(recId, status string) (models.Reception, error) {
	var reception models.Reception

	err := rr.db.QueryRow(updateReceptionStatus,
		status,
		recId,
	).Scan(&reception.Id,
		&reception.DateTime,
		&reception.PVZId,
		&reception.Status,
	)
	if err != nil {
		return reception, dto.ErrDBUpdate
	}

	return reception, nil
}

func (rr *ReceptionRepository) GetReceptionsByPVZIds(pvzIds []string, startDate, endDate *time.Time) ([]models.Reception, error) {
	var rows *sql.Rows
	var err error

	if startDate == nil && endDate == nil {
		rows, err = rr.db.Query(getReceptionsByPVZIds, pq.Array(pvzIds))
	} else {
		start := time.Time{}
		end := time.Now()

		if startDate != nil {
			start = *startDate
		}

		if endDate != nil {
			end = *endDate
		}

		rows, err = rr.db.Query(getReceptionsByPVZIdsFiltered, pq.Array(pvzIds), start, end)
	}

	if err != nil {
		return nil, dto.ErrDBRead
	}
	defer rows.Close()

	receptions := []models.Reception{}

	for rows.Next() {
		var reception models.Reception
		err = rows.Scan(
			&reception.Id,
			&reception.DateTime,
			&reception.PVZId,
			&reception.Status,
		)
		if err != nil {
			return nil, dto.ErrDBRead
		}
		receptions = append(receptions, reception)
	}

	if err = rows.Err(); err != nil {
		return nil, dto.ErrDBRead
	}

	return receptions, nil
}
