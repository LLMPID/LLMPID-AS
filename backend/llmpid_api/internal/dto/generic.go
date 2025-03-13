package dto

type GenericResponse struct {
	Result string `json:"result"`
	Data   string `json:"data"`
}

type Error struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
