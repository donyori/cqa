package cmp

import (
	"reflect"
)

func IsNil(itf interface{}) bool {
	if itf == nil {
		return true
	}
	v := reflect.ValueOf(itf)
	return isNil(v)
}

func IsNilDeeply(itf interface{}) bool {
	if itf == nil {
		return true
	}
	v := reflect.ValueOf(itf)
	if !v.IsValid() {
		return true
	}
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return isNil(v)
}

func isNil(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Ptr:
		fallthrough
	case reflect.Interface:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Chan:
		fallthrough
	case reflect.Func:
		return v.IsNil()
	case reflect.Invalid:
		return true
	default:
		return false
	}
}
