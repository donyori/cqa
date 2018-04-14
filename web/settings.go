package web

import (
	"fmt"

	"github.com/donyori/cqa/common/json"
)

type Settings struct {
	Port uint16 `json:"port"`

	TemplatesPattern string `json:"templates_pattern"`
}

const SettingsFilename string = "../settings/web.json"

var GlobalSettings Settings

func init() {
	// Default values:
	GlobalSettings.Port = 80
	GlobalSettings.TemplatesPattern = "../web/templates/*.tmpl"

	_, err := json.DecodeJsonFromFileIfExist(
		SettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}

func (s *Settings) GetAddr() string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf(":%d", s.Port)
}
