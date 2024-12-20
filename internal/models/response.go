package models

type ErrorResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message"`
}
