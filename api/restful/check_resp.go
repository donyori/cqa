package restful

import (
	"fmt"

	"gopkg.in/resty.v1"
)

func CheckResponse(resp *resty.Response) error {
	if resp == nil {
		return nil
	}
	if respErr := resp.Error(); respErr != nil {
		return fmt.Errorf("%v", respErr)
	}
	if statusCode := resp.StatusCode(); statusCode < 200 || statusCode >= 300 {
		return fmt.Errorf("status: %s", resp.Status())
	}
	return nil
}

func CheckResponseOnAfterResponse(client *resty.Client,
	resp *resty.Response) error {
	return CheckResponse(resp)
}
