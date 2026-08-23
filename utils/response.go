package utils

import (
	"encoding/json"
	"net/http"
)

// APIResponse adalah format response
// standar untuk seluruh API.
type APIResponse struct {

	// Status response.
	// Contoh: success / error
	Status string `json:"status"`

	// Message dari API.
	Message string `json:"message"`

	// Data yang dikirim API.
	Data interface{} `json:"data"`
}

// Success digunakan untuk mengirim
// response ketika request berhasil.
func Success(
	w http.ResponseWriter,
	statusCode int,
	message string,
	data interface{},
) {

	// Response berupa JSON.
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	// Set HTTP status code.
	w.WriteHeader(statusCode)

	// Buat object response.
	response := APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}

	// Convert response menjadi JSON
	// lalu kirim ke client.
	json.NewEncoder(w).Encode(response)
}

// Error digunakan untuk mengirim
// response ketika request gagal.
func Error(
	w http.ResponseWriter,
	statusCode int,
	message string,
) {

	// Response berupa JSON.
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	// Set HTTP status code.
	w.WriteHeader(statusCode)

	// Buat object response error.
	response := APIResponse{
		Status:  "error",
		Message: message,
		Data:    nil,
	}

	// Convert menjadi JSON.
	json.NewEncoder(w).Encode(response)
}
