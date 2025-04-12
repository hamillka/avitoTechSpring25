package repositories

import (
	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
	"github.com/hamillka/avitoTechSpring25/internal/models"
	"github.com/jmoiron/sqlx"
)

type PVZRepository struct {
	db *sqlx.DB
}

const (
	createPVZ             = "INSERT INTO pvzs (city) VALUES ($1) RETURNING id, registration_date, city"
	getPVZById            = "SELECT id, registration_date, city FROM pvzs WHERE id = $1"
	getPVZsWithPagination = "SELECT id, registration_date, city FROM pvzs ORDER BY registration_date DESC LIMIT $1 OFFSET $2"
)

func NewPVZRepository(db *sqlx.DB) *PVZRepository {
	return &PVZRepository{
		db: db,
	}
}

func (pvzr *PVZRepository) CreatePVZ(city string) (models.PVZ, error) {
	var pvz models.PVZ

	err := pvzr.db.QueryRow(createPVZ, city).
		Scan(
			&pvz.Id,
			&pvz.RegistrationDate,
			&pvz.City,
		)
	if err != nil {
		return models.PVZ{}, dto.ErrDBInsert
	}

	return pvz, nil
}

func (pvzr *PVZRepository) GetPVZById(pvzId string) (models.PVZ, error) {
	var pvz models.PVZ
	err := pvzr.db.QueryRow(getPVZById, pvzId).
		Scan(
			&pvz.Id,
			&pvz.RegistrationDate,
			&pvz.City,
		)

	if err != nil {
		return models.PVZ{}, dto.ErrPVZNotFound
	}

	return pvz, nil
}

func (pvzr *PVZRepository) GetPVZsWithPagination(offset, limit int) ([]models.PVZ, error) {
	pvzs := []models.PVZ{}

	rows, err := pvzr.db.Query(getPVZsWithPagination, limit, offset)
	if err != nil {
		return nil, dto.ErrDBRead
	}
	defer rows.Close()

	for rows.Next() {
		var pvz models.PVZ
		err = rows.Scan(
			&pvz.Id,
			&pvz.RegistrationDate,
			&pvz.City,
		)
		if err != nil {
			return nil, dto.ErrDBRead
		}
		pvzs = append(pvzs, pvz)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pvzs, nil
}
