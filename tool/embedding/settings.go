package embedding

import (
	"errors"
	"runtime"

	"github.com/donyori/cqa/common/json"
)

type Settings struct {
	GoroutineNumber       int  `json:"goroutine_number"`
	MinMillisecond        int  `json:"min_millisecond"`
	EnableQuestionBuffer  bool `json:"enable_question_buffer"`
	SkipEmbeddedQuestions bool `json:"skip_embedded_questions"`
	LogStep               int  `json:"log_step"`
}

const SettingsFilename string = "../settings/tool/embedding.json"

var (
	GlobalSettings Settings

	ErrNonPositiveGoroutineNumber error = errors.New(
		"goroutine number is non-positive")
)

func init() {
	// Default values:
	GlobalSettings.GoroutineNumber = runtime.NumCPU() * 16
	GlobalSettings.MinMillisecond = 100
	GlobalSettings.EnableQuestionBuffer = true
	GlobalSettings.SkipEmbeddedQuestions = true
	GlobalSettings.LogStep = 20

	_, err := json.DecodeJsonFromFileIfExist(
		SettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}
