package models

type PVZ struct {
	Id               string
	RegistrationDate string
	City             string
}

type PVZWithReceptions struct {
	PVZ        PVZ
	Receptions []ReceptionWithProducts
}
