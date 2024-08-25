package tools

import (
	"fmt"
	"reflect"
)

func PrettyPrintStruct(dataName string, s interface{}) string {
	res := fmt.Sprintf("\n***** %s *****\n", dataName)

	val := reflect.ValueOf(s)

	// ポインタであれば解参照する
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if val.Field(i).Kind() == reflect.Array && val.Field(i).Type().Elem().Kind() == reflect.Uint8 {
				res += fmt.Sprintf("%s([%d]bytes): %s\n", typ.Field(i).Name, len(val.Field(i).Bytes()), val.Field(i))
			} else {
				res += fmt.Sprintf("%s(%s):  %#v\n", typ.Field(i).Name, val.Field(i).Kind(), val.Field(i))
			}
		}
	default:
		res += fmt.Sprintf("%#v\n", val)
	}

	return res
}
