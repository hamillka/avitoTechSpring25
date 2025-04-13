package models

type PVZ struct {
	Id               string `db:"id"`
	RegistrationDate string `db:"registration_date"`
	City             string `db:"city"`
}

type PVZWithReceptions struct {
	PVZ        PVZ
	Receptions []ReceptionWithProducts
}
