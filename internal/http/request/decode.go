package request

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrEmptyBody   = errors.New("empty body")
	ErrInvalidJSON = errors.New("invalid json")
)

func Decode(r *http.Request, dst any) error {
	if r.Body == nil {
		return ErrEmptyBody
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return ErrInvalidJSON
	}

	return nil
}
