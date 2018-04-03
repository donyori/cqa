package maintain

import (
	"github.com/donyori/cqa/common/json"
)

type EnsureIndexesSettings struct {
	IsBackground bool `json:"is_background"`
}

type MaintainSettings struct {
	EnsureIndexes *EnsureIndexesSettings `json:"ensure_indexes"`
}

const (
	MaintainSettingsFilename string = "settings/maintain.json"
)

var (
	GlobalSettings MaintainSettings
)

func init() {
	// Default values:
	GlobalSettings.EnsureIndexes = &EnsureIndexesSettings{
		IsBackground: true,
	}

	_, err := json.DecodeJsonFromFileIfExist(
		MaintainSettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}
