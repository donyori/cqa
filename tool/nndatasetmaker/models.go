package nndatasetmaker

import (
	"errors"

	"github.com/donyori/cqa/data/model"
)

type IdDatasets struct {
	Train []model.Id `json:"train" bson:"train"`
	Eval  []model.Id `json:"eval" bson:"eval"`
}

type MetaValue struct {
	IdDatasetsMap map[string]IdDatasets `json:"id_datasets" bson:"id_datasets"`
}

type Meta struct {
	model.MetaKey `json:",inline" bson:",inline" cqadm:",inline"`
	Value         *MetaValue `json:"value" bson:"value"`
}

const MetaKey string = "nndatasetmaker"

var ErrCannotGetMeta error = errors.New("cannot get nndatasetmaker meta")

func NewMetaValue() *MetaValue {
	mv := new(MetaValue)
	mv.IdDatasetsMap = make(map[string]IdDatasets)
	return mv
}

func NewMeta() *Meta {
	m := new(Meta)
	m.Key = MetaKey
	return m
}
