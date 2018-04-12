package nlp

import (
	"github.com/donyori/cqa/common/json"
)

type Settings struct {
	HostUrl       string `json:"host_url"`
	EmbeddingPath string `json:"embedding_path"`
}

const SettingsFilename string = "../settings/nlp_api.json"

var GlobalSettings Settings

func init() {
	// Default values:
	GlobalSettings.HostUrl = "http://localhost:5000"
	GlobalSettings.EmbeddingPath = "/embedding"

	_, err := json.DecodeJsonFromFileIfExist(
		SettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}
