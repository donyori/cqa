package stackexchange

import (
	"time"

	"github.com/donyori/cqa/common/json"
)

type Settings struct {
	// API info.
	HostUrl    string `json:"host_url"`
	ApiVersion string `json:"api_version"`

	// API path.
	QuestionsPath string `json:"questions_path"`

	// APP info.
	ClientId    int    `json:"client_id"`
	AccessToken string `json:"access_token"`
	Key         string `json:"key"`

	// Common query params.
	Filter string `json:"filter"`
	Site   string `json:"site"`

	// Paging.
	StartPage       int `json:"start_page"`
	DefaultPageSize int `json:"default_page_size"`
	MaxPageSize     int `json:"max_page_size"`

	// Rate limiting.
	BackoffUnit time.Duration `json:"backoff_unit"`
}

const SettingsFilename string = "../settings/stack_exchange_api.json"

var GlobalSettings Settings

func init() {
	// Default values:
	GlobalSettings.HostUrl = "https://api.stackexchange.com"
	GlobalSettings.ApiVersion = "2.2"

	GlobalSettings.QuestionsPath = "/questions"

	GlobalSettings.ClientId = 11473
	GlobalSettings.AccessToken = `u(3incSNPFr9j7yNKlaDqg))`
	GlobalSettings.Key = `omh9)pqtJwhAiAYQZmLkKA((`

	GlobalSettings.Filter =
		`!2A0mshvsMdcL9mpcSwgbRIOW0.ErSpw*FW__y3sXqHOTysHny(Ly4D`
	GlobalSettings.Site = "stackoverflow"

	GlobalSettings.StartPage = 1
	GlobalSettings.DefaultPageSize = 30
	GlobalSettings.MaxPageSize = 100

	GlobalSettings.BackoffUnit = time.Second

	_, err := json.DecodeJsonFromFileIfExist(
		SettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}
