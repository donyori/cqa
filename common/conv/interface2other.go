package conv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

var (
	ErrCannotToInt64     error = errors.New("interface cannot convert to int64")
	ErrCannotToTimestamp error = errors.New(
		"interface cannot convert to timestamp")
)

func InterfaceToInt64(itf interface{}) (i64 int64, err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			i64 = 0
			e, ok := panicErr.(error)
			if ok {
				err = e
			} else {
				err = fmt.Errorf("%v", panicErr)
			}
		}
	}()
	if itf == nil {
		return 0, ErrCannotToInt64
	}
	v := reflect.ValueOf(itf)
	return toInt64(v)
}

func InterfaceToTimestamp(itf interface{}) (timestamp int64, err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			timestamp = 0
			e, ok := panicErr.(error)
			if ok {
				err = e
			} else {
				err = fmt.Errorf("%v", panicErr)
			}
		}
	}()
	if itf == nil {
		return 0, ErrCannotToTimestamp
	}
	v := reflect.ValueOf(itf)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return 0, ErrCannotToTimestamp
	}
	if v.Kind() == reflect.Struct {
		if v.Type() != reflect.TypeOf(time.Time{}) {
			return 0, ErrCannotToTimestamp
		}
		i := v.Interface()
		t, ok := i.(time.Time)
		if !ok {
			return 0, ErrCannotToTimestamp
		}
		return t.Unix(), nil
	} else if v.Kind() == reflect.Invalid {
		return 0, ErrCannotToTimestamp
	} else {
		i64, err := toInt64(v)
		if err == nil {
			return i64, nil
		}
		return 0, ErrCannotToTimestamp
	}
}

func toInt64(v reflect.Value) (i64 int64, err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			i64 = 0
			e, ok := panicErr.(error)
			if ok {
				err = e
			} else {
				err = fmt.Errorf("%v", panicErr)
			}
		}
	}()
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return 0, ErrCannotToInt64
	}
	switch v.Kind() {
	case reflect.String:
		s := v.String()
		i64, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			return 0, ErrCannotToInt64
		}
		return i64, nil
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		return v.Int(), nil
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		return int64(v.Uint()), nil
	default:
		return 0, ErrCannotToInt64
	}
}
