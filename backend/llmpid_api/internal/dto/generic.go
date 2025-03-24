package dto

type GenericResponse struct {
	Status  string `json:"result"`
	Message string `json:"data"`
}
