package response

import "net/http"

func NewResponse(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}
