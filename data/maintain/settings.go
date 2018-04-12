package maintain

import (
	"github.com/donyori/cqa/common/json"
)

type EnsureIndexesSettings struct {
	IsBackground bool `json:"is_background"`
}

type Settings struct {
	EnsureIndexes *EnsureIndexesSettings `json:"ensure_indexes"`
}

const SettingsFilename string = "../settings/maintain.json"

var GlobalSettings Settings

func init() {
	// Default values:
	GlobalSettings.EnsureIndexes = &EnsureIndexesSettings{
		IsBackground: true,
	}

	_, err := json.DecodeJsonFromFileIfExist(
		SettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}
