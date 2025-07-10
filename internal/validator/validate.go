package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func Validate(ptr any) error {
	if ptr == nil {
		return errors.New("validate: nil value")
	}

	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("validate: expect non-nil pointer to struct")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("validate: expect pointer to struct")
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		tag := fld.Tag.Get("validate")
		name := fld.Tag.Get("json")
		name = strings.TrimSuffix(name, ",omitempty")
		if tag == "" {
			continue
		}
		value := v.Field(i).Interface()
		for _, rule := range strings.Split(tag, ",") {
			if rule == "" {
				continue
			}
			if err := applyRule(name, rule, value); err != nil {
				return err
			}
		}
	}

	return nil
}

func applyRule(name, rule string, val any) error {
	switch {
	case rule == "required":
		if isZero(val) {
			return fmt.Errorf("%s is required", name)
		}
	case rule == "string":
		if reflect.TypeOf(val).Kind() != reflect.String {
			return fmt.Errorf("%s must be a string", name)
		}
	}

	return nil
}

func isZero(v any) bool {
	return reflect.ValueOf(v).IsZero()
}
