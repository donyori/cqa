package tag

import (
	"reflect"
	"strings"
)

type dfsState struct {
	Value reflect.Value
	Index int
}

const (
	TagInline string = "inline"

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
	tagValue = strings.TrimSpace(tagValue)
	field = nil
	ok = false
	f := func(tag string, value *reflect.Value) bool {
		if tag == tagValue {
			field = value
			ok = true
			return true
		}
		return false
	}
	dfs(model, tagKey, f)
	return field, ok
}

func GetDataModelFieldValueByTag(model interface{},
	tagKey string, tagValue string) (
	value interface{}, ok bool) {
	field, ok := GetDataModelFieldByTag(model, tagKey, tagValue)
	if !ok {
		return nil, false
	}
	if field == nil || !field.IsValid() {
		return nil, true
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
	tagValue = strings.TrimSpace(tagValue)
	fields = nil
	f := func(tag string, value *reflect.Value) bool {
		if tag == tagValue {
			fields = append(fields, value)
		}
		return false
	}
	dfs(model, tagKey, f)
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
	fieldsMap = nil
	f := func(tag string, value *reflect.Value) bool {
		if fieldsMap == nil {
			fieldsMap = make(map[string][]*reflect.Value)
		}
		fieldsMap[tag] = append(fieldsMap[tag], value)
		return false
	}
	dfs(model, tagKey, f)
	return fieldsMap
}

func dfs(model interface{}, tagKey string,
	f func(string, *reflect.Value) bool) {
	if model == nil {
		return
	}
	v := reflect.ValueOf(model)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}
	var dfsStack []dfsState
	dfsStack = append(dfsStack, dfsState{Value: v, Index: 0})
	top := 0
	// DFS
	for top >= 0 {
		v = dfsStack[top].Value
		t := v.Type()
		numField := t.NumField()
		hasInline := false
		for i := dfsStack[top].Index; i < numField; i++ {
			fieldType := t.Field(i)
			tagsStr, hasTag := fieldType.Tag.Lookup(tagKey)
			if !hasTag {
				continue
			}
			tags := strings.Split(tagsStr, ",")
			for _, tag := range tags {
				tag = strings.TrimSpace(tag)
				if tag != TagInline {
					fv := v.Field(i)
					isExit := f(tag, &fv)
					if isExit {
						return
					}
				} else {
					hasInline = true
				}
			}
			if hasInline {
				fv := v.Field(i)
				for fv.Kind() == reflect.Ptr {
					fv = fv.Elem()
				}
				if fv.Kind() != reflect.Struct {
					hasInline = false
					continue
				}
				dfsStack[top].Index = i + 1
				top++
				if top >= len(dfsStack) {
					dfsStack = append(dfsStack, dfsState{
						Value: fv,
						Index: 0,
					})
				} else {
					dfsStack[top].Value = fv
					dfsStack[top].Index = 0
				}
				break
			}
		}
		if !hasInline {
			top--
		}
	}
}
