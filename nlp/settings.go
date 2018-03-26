package nlp

import (
	"errors"

	"github.com/donyori/cqa/common/json"
)

type NlpSettings struct {
	EmbeddingMethod string `json:"embedding_method"`
}

const (
	NlpMethodUseRestfulApi string = "use_restful_api"

	NlpSettingsFilename string = "settings/nlp.json"
)

var (
	GlobalSettings NlpSettings

	ErrUnknownNlpMethod error = errors.New("unknown method for NLP")
)

func init() {
	// Default values:
	GlobalSettings.EmbeddingMethod = NlpMethodUseRestfulApi

	_, err := json.DecodeJsonFromFileIfExist(NlpSettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}
