package models

type Reception struct {
	Id       string
	DateTime string
	PVZId    string
	Status   string
}

type ReceptionWithProducts struct {
	Reception Reception
	Products  []Product
}

const (
	CLOSE      = "close"
	INPROGRESS = "in_progress"
)
