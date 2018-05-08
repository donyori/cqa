package crawler

import (
	"github.com/donyori/cqa/common/json"
)

type Settings struct {
	CrawlTags []string `json:"crawl_tags"`
	LogStep   int      `json:"log_step"`
}

const SettingsFilename string = "../settings/crawler.json"

var GlobalSettings Settings

func init() {
	// Default values:
	GlobalSettings.CrawlTags = []string{"c", "c++"}
	GlobalSettings.LogStep = 10

	_, err := json.DecodeJsonFromFileIfExist(
		SettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}
