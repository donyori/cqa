package nlp

import (
	"github.com/donyori/cqa/common/json"
)

type NlpApiSettings struct {
	HostUrl       string `json:"host_url"`
	EmbeddingPath string `json:"embedding_path"`
}

const (
	NlpApiSettingsFilename string = "settings/nlp_api.json"
)

var (
	GlobalSettings NlpApiSettings
)

func init() {
	// Default values:
	GlobalSettings.HostUrl = "http://localhost:5000"
	GlobalSettings.EmbeddingPath = "/embedding"

	_, err := json.DecodeJsonFromFileIfExist(NlpApiSettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}
