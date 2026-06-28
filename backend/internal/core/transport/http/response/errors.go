package http_response

type ErrorResponse struct {
	Error   string `json:"error"   example:"user with id='42': not found"`
	Message string `json:"message" example:"failed to get user"`
}
