package nlp

import (
	"sync"

	"gopkg.in/resty.v1"

	"github.com/donyori/cqa/api/restful"
)

var (
	nlpGetClient *resty.Client

	initNlpGetClientOnce sync.Once
)

func GetNlpGetClient() *resty.Client {
	initNlpGetClientOnce.Do(initNlpGetClient)
	return nlpGetClient
}

func initNlpGetClient() {
	nlpGetClient = resty.New().
		SetHostURL(GlobalSettings.HostUrl).
		SetHeader("Accept", "application/json").
		OnAfterResponse(restful.CheckResponseOnAfterResponse)
}
