package embedding

import (
	"runtime"

	"github.com/donyori/cqa/common/json"
)

type EmbeddingSettings struct {
	GoroutineNumber int `json:"goroutine_number"`
	MinMillisecond  int `json:"min_millisecond"`
}

const EmbeddingSettingsFilename string = "settings/embedding.json"

var GlobalSettings EmbeddingSettings

func init() {
	// Default values:
	GlobalSettings.GoroutineNumber = runtime.NumCPU() * 16
	GlobalSettings.MinMillisecond = 250

	_, err := json.DecodeJsonFromFileIfExist(
		EmbeddingSettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}
