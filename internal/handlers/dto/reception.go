package dto

// ReceptionDto model info
// @Description Информация о приемке
type ReceptionDto struct {
	Id       string `json:"id"`       // Идентификатор приемки
	DateTime string `json:"dateTime"` // Дата и время создания приемки
	PVZId    string `json:"pvzId"`    // Идентификатор ПВЗ
	Status   string `json:"status"`   // Статус приемки
}

// ReceptionWithProducts model info
// @Description Информация о приемке и товарах в ней
type ReceptionWithProductsDto struct {
	Reception ReceptionDto `json:"reception"` // Информация о приемке
	Products  []ProductDto `json:"products"`  // Информация о всех товарах в приемке
}

// CreateReceptionRequestDto model info
// @Description Информация о приемке при ее создании
type CreateReceptionRequestDto struct {
	PVZId string `json:"pvzId"` // Идентификатор ПВЗ
}

// CreateReceptionResponseDto model info
// @Description Информация о приемке при ее создании
type CreateReceptionResponseDto struct {
	Id       string `json:"id"`       // Идентификатор приемки
	DateTime string `json:"dateTime"` // Дата и время приемки
	PVZId    string `json:"pvzId"`    // Идентификатор ПВЗ
	Status   string `json:"status"`   // Статус приемки
}

// CloseReceptionResponseDto model info
// @Description Информация о приемке при ее закрытии
type CloseReceptionResponseDto struct {
	Id       string `json:"id"`       // Идентификатор приемки
	DateTime string `json:"dateTime"` // Дата и время приемки
	PVZId    string `json:"pvzId"`    // Идентификатор ПВЗ
	Status   string `json:"status"`   // Статус приемки
}
