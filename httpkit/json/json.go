package json

import (
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
)

type err struct {
	Err string `json:"error"`
}

func NewError(Err string) *err {
	return &err{
		Err: Err,
	}
}

func (e *err) Error() string {
	return fmt.Sprintf("{\n\"error\": \"%s\"}\n", e.Err)
}

func OK(w http.ResponseWriter, data any) error {
	return Write(w, http.StatusOK, data)
}

func Unauthorized(w http.ResponseWriter, data any) error {
	return Write(w, http.StatusUnauthorized, data)
}

func BadRequest(w http.ResponseWriter, data any) error {
	return Write(w, http.StatusBadRequest, data)
}

func Write(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err = w.Write(json); err != nil {
		return err
	}

	return nil
}

func InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(500)
	_, _ = w.Write([]byte(NewError("internal server error").Error()))
}
