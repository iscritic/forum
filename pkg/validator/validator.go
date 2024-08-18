package validator

import (
	"reflect"
	"regexp"
)

func IsEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Slice,
		reflect.Map, reflect.Ptr, reflect.Interface:
		return v.IsNil()
	}
	return false
}

func IsLengthValid(str string, min, max int) bool {
	length := len(str)
	return length >= min && length <= max
}

func IsInSet(value interface{}, set ...interface{}) bool {
	for _, item := range set {
		if value == item {
			return true
		}
	}
	return false
}

func IsEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)

}
