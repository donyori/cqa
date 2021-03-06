package qm

import (
	"runtime"
	"time"

	"github.com/donyori/cqa/common/json"
)

type ExitMode int8

type Settings struct {
	DefaultTopNumber            int      `json:"default_top_number"`
	DefaultTimeLimitMillisecond int64    `json:"default_time_limit_millisecond"`
	DefaultExitMode             ExitMode `json:"default_exit_mode"`
	DefaultMatcherNumber        int      `json:"matcher_number"`

	RequestBufferSize uint32 `json:"request_buffer_size"`
	ErrorBufferSize   uint32 `json:"error_buffer_size"`
	InputBufferSize   uint32 `json:"input_buffer_size"`
	OutputBufferSize  uint32 `json:"output_buffer_size"`

	EnableQuestionVectorBuffer bool `json:"enable_question_vector_buffer"`
}

const (
	ExitModeGracefully ExitMode = iota
	ExitModeImmediately
	ExitModeForcedly
)

const SettingsFilename string = "../settings/qm.json"

var (
	exitModeStrings = [...]string{
		"Exit Gracefully",
		"Exit Immediately",
		"Exit Forcedly",
	}

	GlobalSettings Settings
)

func init() {
	// Default values:
	GlobalSettings.DefaultTopNumber = 5
	GlobalSettings.DefaultTimeLimitMillisecond = 0
	GlobalSettings.DefaultExitMode = ExitModeGracefully
	GlobalSettings.DefaultMatcherNumber = runtime.NumCPU()

	GlobalSettings.RequestBufferSize = 5
	GlobalSettings.ErrorBufferSize = uint32(
		(GlobalSettings.DefaultMatcherNumber + 1) * 2)
	GlobalSettings.InputBufferSize = 5
	GlobalSettings.OutputBufferSize = uint32(
		GlobalSettings.DefaultMatcherNumber)

	GlobalSettings.EnableQuestionVectorBuffer = true

	_, err := json.DecodeJsonFromFileIfExist(
		SettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}

func (em ExitMode) IsValid() bool {
	return em >= ExitModeGracefully && em <= ExitModeForcedly
}

func (em ExitMode) String() string {
	if !em.IsValid() {
		return "unknown"
	}
	return exitModeStrings[em]
}

func (s *Settings) GetDefaultTimeLimit() time.Duration {
	if s == nil {
		return 0
	}
	var tl time.Duration
	if s.DefaultTimeLimitMillisecond > 0 {
		tl = time.Millisecond * time.Duration(s.DefaultTimeLimitMillisecond)
	} else {
		tl = time.Duration(s.DefaultTimeLimitMillisecond)
	}
	return tl
}
