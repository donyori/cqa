package nlp

import (
	"errors"

	"github.com/donyori/cqa/common/json"
)

type NlpSource int8

type Settings struct {
	EmbeddingSource           NlpSource `json:"embedding_source"`
	EmbeddingWithTokensSource NlpSource `json:"embedding_with_tokens_source"`
}

const (
	NlpSourceRestfulApi NlpSource = iota
)

const SettingsFilename string = "../settings/nlp.json"

var (
	GlobalSettings Settings

	ErrInvalidNlpSource error = errors.New("NLP source is invalid")

	nlpSourceStrings = [...]string{
		"RESTful API",
	}
)

func init() {
	// Default values:
	GlobalSettings.EmbeddingSource = NlpSourceRestfulApi
	GlobalSettings.EmbeddingWithTokensSource = NlpSourceRestfulApi

	_, err := json.DecodeJsonFromFileIfExist(
		SettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}

func (ns NlpSource) IsValid() bool {
	return ns == NlpSourceRestfulApi
}

func (ns NlpSource) String() string {
	if !ns.IsValid() {
		return "Unknown"
	}
	return nlpSourceStrings[ns]
}
