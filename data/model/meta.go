package model

import (
	"time"
)

type MetaKey struct {
	Key string `json:"key" bson:"_id" cqadm:"id"`
}

type MetaInterface struct {
	MetaKey `json:",inline" bson:",inline" cqadm:",inline"`
	Value   interface{} `json:"value" bson:"value"`
}

type MetaBool struct {
	MetaKey `json:",inline" bson:",inline" cqadm:",inline"`
	Value   bool `json:"value" bson:"value"`
}

type MetaInt64 struct {
	MetaKey `json:",inline" bson:",inline" cqadm:",inline"`
	Value   int64 `json:"value" bson:"value"`
}

type MetaUint64 struct {
	MetaKey `json:",inline" bson:",inline" cqadm:",inline"`
	Value   uint64 `json:"value" bson:"value"`
}

type MetaFloat64 struct {
	MetaKey `json:",inline" bson:",inline" cqadm:",inline"`
	Value   float64 `json:"value" bson:"value"`
}

type MetaComplex128 struct {
	MetaKey `json:",inline" bson:",inline" cqadm:",inline"`
	Value   complex128 `json:"value" bson:"value"`
}

type MetaString struct {
	MetaKey `json:",inline" bson:",inline" cqadm:",inline"`
	Value   string `json:"value" bson:"value"`
}

type MetaTime struct {
	MetaKey `json:",inline" bson:",inline" cqadm:",inline"`
	Value   *time.Time `json:"value" bson:"value"`
}

func NewMetaInterface() *MetaInterface {
	return new(MetaInterface)
}

func NewMetaBool() *MetaBool {
	return new(MetaBool)
}

func NewMetaInt64() *MetaInt64 {
	return new(MetaInt64)
}

func NewMetaUint64() *MetaUint64 {
	return new(MetaUint64)
}

func NewMetaFloat64() *MetaFloat64 {
	return new(MetaFloat64)
}

func NewMetaComplex128() *MetaComplex128 {
	return new(MetaComplex128)
}

func NewMetaString() *MetaString {
	return new(MetaString)
}

func NewMetaTime() *MetaTime {
	return new(MetaTime)
}
