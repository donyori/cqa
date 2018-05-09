package tagcount

import (
	"errors"
	"runtime"

	"github.com/donyori/cqa/common/json"
)

type Settings struct {
	GoroutineNumber int    `json:"goroutine_number"`
	OutputFilename  string `json:"output_filename"`
	LogStep         int    `json:"log_step"`
}

const SettingsFilename string = "../settings/tool/tagcount.json"

var (
	GlobalSettings Settings

	ErrNonPositiveGoroutineNumber error = errors.New(
		"goroutine number is non-positive")
)

func init() {
	// Default values:
	GlobalSettings.GoroutineNumber = runtime.NumCPU()
	GlobalSettings.OutputFilename = "../out/tagcount/result.csv"
	GlobalSettings.LogStep = 1000

	_, err := json.DecodeJsonFromFileIfExist(
		SettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}
