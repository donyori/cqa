package tag

import (
	"reflect"
)

const (
	DataModelTagKey string = "cqadm"
	DataModelTagId  string = "id"

	MongoDbTagKey string = "bson"
	MongoDbTagId  string = "_id"
)

func GetDataModelFieldByTag(model interface{},
	tagKey string, tagValue string) (
	field *reflect.Value, ok bool) {
	if model == nil {
		return nil, false
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			field = nil
			ok = false
		}
	}()
	v := reflect.ValueOf(model)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, false
	}
	t := v.Type()
	numField := t.NumField()
	var fieldType reflect.StructField
	var tag string
	var hasTag bool
	for i := 0; i < numField; i++ {
		fieldType = t.Field(i)
		tag, hasTag = fieldType.Tag.Lookup(tagKey)
		if hasTag && tagValue == tag {
			fieldValue := v.Field(i)
			return &fieldValue, true
		}
	}
	return nil, false
}

func GetDataModelFieldValueByTag(model interface{},
	tagKey string, tagValue string) (
	value interface{}, ok bool) {
	field, ok := GetDataModelFieldByTag(model, tagKey, tagValue)
	if !ok {
		return nil, false
	}
	return field.Interface(), true
}

func GetDataModelFieldsByTag(model interface{},
	tagKey string, tagValue string) (
	fields []*reflect.Value) {
	if model == nil {
		return nil
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			fields = nil
		}
	}()
	v := reflect.ValueOf(model)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	t := v.Type()
	numField := t.NumField()
	var fieldType reflect.StructField
	var tag string
	var hasTag bool
	fields = nil
	for i := 0; i < numField; i++ {
		fieldType = t.Field(i)
		tag, hasTag = fieldType.Tag.Lookup(tagKey)
		if hasTag && tagValue == tag {
			fieldValue := v.Field(i)
			fields = append(fields, &fieldValue)
		}
	}
	return fields
}

func GetDataModelFieldsGroupByTag(model interface{}, tagKey string) (
	fieldsMap map[string][]*reflect.Value) {
	if model == nil {
		return nil
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			fieldsMap = nil
		}
	}()
	v := reflect.ValueOf(model)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	t := v.Type()
	numField := t.NumField()
	var fieldType reflect.StructField
	var tagValue string
	var hasTag bool
	fieldsMap = nil
	for i := 0; i < numField; i++ {
		fieldType = t.Field(i)
		tagValue, hasTag = fieldType.Tag.Lookup(tagKey)
		if hasTag {
			fieldValue := v.Field(i)
			if fieldsMap == nil {
				fieldsMap = make(map[string][]*reflect.Value)
			}
			fieldsMap[tagValue] = append(fieldsMap[tagValue], &fieldValue)
		}
	}
	return fieldsMap
}
