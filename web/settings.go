package web

import (
	"fmt"
	"net/http"

	"github.com/donyori/cqa/common/json"
)

type Settings struct {
	Port uint16 `json:"port"`

	StaticResourcesRoot http.Dir `json:"static_resources_root"`
	TemplatesPattern    string   `json:"templates_pattern"`
}

const SettingsFilename string = "../settings/web.json"

var GlobalSettings Settings

func init() {
	// Default values:
	GlobalSettings.Port = 80
	GlobalSettings.StaticResourcesRoot = http.Dir("../web/static_resources")
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
