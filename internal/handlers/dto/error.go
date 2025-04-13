package dto

import (
	goErrors "errors"
)

var (
	ErrPVZNotFound            = goErrors.New("no such PVZ")
	ErrNoActiveReception      = goErrors.New("no active reception")
	ErrPVZReceptionIsClosed   = goErrors.New("reception is closed")
	ErrNoProductsInReception  = goErrors.New("no products in reception")
	ErrPVZAlreadyHasReception = goErrors.New("PVZ already has active reception")
	ErrUserAlreadyExists      = goErrors.New("user already exists")
	ErrInvalidCredentials     = goErrors.New("user login invalid credentials")
	ErrDBInsert               = goErrors.New("failed to insert into DB")
	ErrDBRead                 = goErrors.New("failed to read from DB")
	ErrDBUpdate               = goErrors.New("failer to update in DB")
)

// ErrorDto model info
// @Description Информация об ошибке (DTO)
type ErrorDto struct {
	Message string `json:"message"` // Текст ошибки
}
