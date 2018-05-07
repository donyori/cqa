package stackexchange

import (
	"sync"

	"gopkg.in/resty.v1"

	"github.com/donyori/cqa/api/restful"
)

var (
	seGetClient *resty.Client

	initSeGetClientOnce sync.Once
)

func GetSeGetClient() *resty.Client {
	initSeGetClient()
	return seGetClient
}

func initSeGetClient() {
	initSeGetClientOnce.Do(func() {
		seGetClient = resty.New().
			SetHostURL(GlobalSettings.HostUrl).
			SetHeader("Accept", "application/json").
			OnAfterResponse(restful.CheckResponseOnAfterResponse).
			SetQueryParams(map[string]string{
				"access_token": GlobalSettings.AccessToken,
				"key":          GlobalSettings.Key,
				"site":         GlobalSettings.Site,
				"filter":       GlobalSettings.Filter,
			})
	})
}
