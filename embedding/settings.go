package embedding

import (
	"errors"
	"runtime"

	"github.com/donyori/cqa/common/json"
)

type Settings struct {
	GoroutineNumber int `json:"goroutine_number"`
	MinMillisecond  int `json:"min_millisecond"`
}

const SettingsFilename string = "settings/embedding.json"

var (
	GlobalSettings Settings

	ErrNonPositiveGoroutineNumber error = errors.New(
		"goroutine number is non-positive")
)

func init() {
	// Default values:
	GlobalSettings.GoroutineNumber = runtime.NumCPU() * 16
	GlobalSettings.MinMillisecond = 250

	_, err := json.DecodeJsonFromFileIfExist(
		SettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}
