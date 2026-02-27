package utils

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/matheusgosk8/book-me-server/internal/models"
)

func errorHandler(res http.ResponseWriter, message string, code int) {

	response := models.ServerStatus{
		Code:    code,
		Message: message,
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	json.NewEncoder(res).Encode(response)

}

var RequestErrorHandler = func(w http.ResponseWriter, err error) {
	errorHandler(w, err.Error(), http.StatusBadRequest)
}
var InternalErrorHandler = func(w http.ResponseWriter, err error) {
	errorHandler(w, "Internal error!", http.StatusInternalServerError)
}

func BodyParser[T any](r *http.Request) (*T, error) {
	data := new(T)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&data); err != nil {
		log.Error("Error decoding request body: ", err)
		return nil, err
	}
	return data, nil
}

func ServerResponse[T any](res http.ResponseWriter, data T) {
	res.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(res).Encode(data)
	if err != nil {
		InternalErrorHandler(res, err)
	}
}

func GenerateID() string {
	return uuid.New().String()
}
