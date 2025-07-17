package utils

import "reflect"

func ApplyUpdates(model any) (sets []string, args []any) {
	v := reflect.ValueOf(model)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		dbTag := t.Field(i).Tag.Get("db")
		arg := v.Field(i).Interface()

		if !reflect.ValueOf(arg).IsZero() {
			sets = append(sets, dbTag+" = ?")
			args = append(args, arg)
		}
	}

	return
}
