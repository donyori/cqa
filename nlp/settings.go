package nlp

import (
	"errors"

	"github.com/donyori/cqa/common/json"
)

type NlpSource int8

type NlpSettings struct {
	EmbeddingSource NlpSource `json:"embedding_source"`
}

const (
	NlpSourceRestfulApi NlpSource = iota
)

const NlpSettingsFilename string = "settings/nlp.json"

var (
	GlobalSettings NlpSettings

	ErrInvalidNlpSource error = errors.New("NLP source is invalid")

	nlpSourceStrings = [...]string{
		"RESTful API",
	}
)

func init() {
	// Default values:
	GlobalSettings.EmbeddingSource = NlpSourceRestfulApi

	_, err := json.DecodeJsonFromFileIfExist(
		NlpSettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}

func (ns NlpSource) String() string {
	if !ns.IsValid() {
		return "Unknown"
	}
	return nlpSourceStrings[ns]
}

func (ns NlpSource) IsValid() bool {
	return ns == NlpSourceRestfulApi
}
