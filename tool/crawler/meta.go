package crawler

import (
	"time"

	"github.com/donyori/cqa/data/model"
)

type MetaValue struct {
	LastCrawlTime     time.Time            `json:"last_crawl_time" bson:"last_crawl_time"`
	LastActivityDates map[string]time.Time `json:"last_activity_dates" bson:"last_activity_dates"`
}

type Meta struct {
	model.MetaKey `json:",inline" bson:",inline" cqadm:",inline"`
	Value         *MetaValue `json:"value" bson:"value"`
}

const MetaKey string = "crawler"

func NewMetaValue() *MetaValue {
	return new(MetaValue)
}

func NewMeta() *Meta {
	m := new(Meta)
	m.Key = MetaKey
	return m
}
