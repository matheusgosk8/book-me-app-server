package utils

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"strconv"

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
	// decoder.DisallowUnknownFields()
	defer r.Body.Close()

	if err := decoder.Decode(data); err != nil {
		log.Error("Error decoding request body: ", err)
		return nil, err
	}
	return data, nil
}

func ServerResponse[T any](res http.ResponseWriter, data T, meta ...models.Meta) {
	ServerSuccess(res, http.StatusOK, "OK", data, meta...)
}

// Standard response format for success
func ServerSuccess[T any](res http.ResponseWriter, statusCode int, message string, data T, meta ...models.Meta) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)

	payload := map[string]any{
		"statusCode": statusCode,
		"message":    message,
	}

	// Include data only when it's not a zero value (allows exclusive responses like deletes using nil)
	// Use reflection safely to detect zero/nil values for generic T.
	includeData := false
	var dataAny any = any(data)
	if dataAny != nil {
		// reflect.ValueOf on non-nil
		// Note: IsZero will return true for zero values; we include data only if not zero.
		v := reflect.ValueOf(dataAny)
		if v.IsValid() && !v.IsZero() {
			includeData = true
		}
	}
	if includeData {
		payload["data"] = data
	}

	if len(meta) > 0 {
		payload["meta"] = meta[0]
	}

	if err := json.NewEncoder(res).Encode(payload); err != nil {
		log.Error("Error encoding success response: ", err)
		InternalErrorHandler(res, err)
	}
}

// Standard response format for errors. Returns provided error message and status code.
func ServerError(res http.ResponseWriter, statusCode int, err error) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)

	payload := map[string]any{
		"statusCode": statusCode,
		"message":    err.Error(),
	}

	if e := json.NewEncoder(res).Encode(payload); e != nil {
		log.Error("Error encoding error response: ", e)
	}
}

func GenerateID() string {
	return uuid.New().String()
}

func QueryInt(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)

	if value == "" {
		return defaultValue
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return result
}
