package titlelen

import (
	"errors"
	"runtime"

	"github.com/donyori/cqa/common/json"
)

type Settings struct {
	GoroutineNumber int    `json:"goroutine_number"`
	OutputDirectory string `json:"output_directory"`
	LogStep         int    `json:"log_step"`
}

const SettingsFilename string = "../settings/tool/titlelen.json"

var (
	GlobalSettings Settings

	ErrNonPositiveGoroutineNumber error = errors.New(
		"goroutine number is non-positive")
)

func init() {
	// Default values:
	GlobalSettings.GoroutineNumber = runtime.NumCPU()
	GlobalSettings.OutputDirectory = "../out/titlelen"
	GlobalSettings.LogStep = 1000

	_, err := json.DecodeJsonFromFileIfExist(
		SettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}
