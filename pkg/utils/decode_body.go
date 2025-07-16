package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func DecodeBody(r *http.Request, dto any) error {
	dec := json.NewDecoder(r.Body)
	//dec.DisallowUnknownFields()

	if err := dec.Decode(dto); err != nil {
		var ute *json.UnmarshalTypeError
		if errors.As(err, &ute) {
			msg := fmt.Sprintf("%s must be %s", strings.ToLower(ute.Field), ute.Type.Kind())
			return errors.New(msg)
		}
		if errors.Is(err, io.EOF) {
			return errors.New("request body must not be empty")
		}

		return err
	}

	return nil
}
