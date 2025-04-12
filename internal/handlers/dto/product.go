package dto

const (
	ProductTypeElectronics = "электроника"
	ProductTypeClothes     = "одежда"
	ProductTypeShoes       = "обувь"
)

// ProductDto model info
// @Description Информация о товаре
type ProductDto struct {
	Id          string `json:"id"`          // Идентификатор
	DateTime    string `json:"dateTime"`    // Дата и время
	Type        string `json:"type"`        // Тип товара
	ReceptionId string `json:"receptionId"` // Идентификатор приемки
}

// AddProductRequestDto model info
// @Description Информация запроса о товаре при его добавлении
type AddProductRequestDto struct {
	Type  string `json:"type"`  // Тип товара
	PVZId string `json:"pvzId"` // Идентификатор ПВЗ, на который добавляется товар
}

// AddProductResponseDto model info
// @Description Информация о товаре при его добавлении
type AddProductResponseDto struct {
	Id          string `json:"id"`          // Идентификатор
	DateTime    string `json:"dateTime"`    // Дата и время
	Type        string `json:"type"`        // Тип товара
	ReceptionId string `json:"receptionId"` // Идентификатор приемки
}
