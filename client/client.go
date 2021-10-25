package client

import (
	"fmt"
	"hubs-cms-go/logger"

	"github.com/go-resty/resty/v2"
)

// RestyClient the resty http client
var RestyClient = resty.New()

// NewHTTPRequest create http request
func NewHTTPRequest() *resty.Request {
	return RestyClient.R()
}

// Setup setup http client
func Setup() {
	RestyClient.OnAfterResponse(responseLogger)
}

func responseLogger(c *resty.Client, resp *resty.Response) error {
	method := fmt.Sprintf("[Method] %s", resp.Request.Method)
	url := fmt.Sprintf("[URL] %s", resp.Request.URL)
	reqBody := fmt.Sprintf("[Request Body] %v", resp.Request.Body)
	status := fmt.Sprintf("[Status] %d", resp.StatusCode())
	duration := fmt.Sprintf("[Duration] %v", resp.Time())
	respBody := fmt.Sprintf("[Response Body] %s", resp.String())

	log := fmt.Sprintf("%s %s %s %s %s %s", method, url, reqBody, status, duration, respBody)
	logger.Debug.Println(log)
	return nil
}
