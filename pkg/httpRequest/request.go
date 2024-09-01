package httpRequest

import (
	"net/http"
	"time"
)

var HTTPClient = &http.Client{
	Timeout: 5 * time.Second,
}

func HTTPDo(req *http.Request) (*http.Response, error) {

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
