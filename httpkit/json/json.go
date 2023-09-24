package json

import (
	"errors"
	"fmt"
	"io"
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
	return fmt.Sprintf("{\n\"error\": \"%s\"\n}\n", e.Err)
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

func UnprocessableEntity(w http.ResponseWriter, data any) error {
	return Write(w, http.StatusUnprocessableEntity, data)
}

func Conflict(w http.ResponseWriter, data any) error {
	return Write(w, http.StatusConflict, data)
}

func InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(500)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(NewError("internal server error").Error()))
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

func Read(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		var (
			syntaxErr           *json.SyntaxError
			unmarshalTypeErr    *json.UnmarshalTypeError
			invalidUnmarshalErr *json.InvalidUnmarshalError
		)

		switch {
		case errors.As(err, &syntaxErr):
			return fmt.Errorf("badly formed JSON (at character %d)", syntaxErr.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("request closed")

		case errors.As(err, &unmarshalTypeErr):
			if unmarshalTypeErr.Field != "" {
				return fmt.Errorf("incorrect type (at field %s)", unmarshalTypeErr.Field)
			}

			return fmt.Errorf("incorrect type (at character %d)", unmarshalTypeErr.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("empty JSON")

		case errors.As(err, &invalidUnmarshalErr):
			panic("destination must be non nil pointer")
		default:
			return err
		}
	}

	return nil
}
