package models

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
)

//ClientRequest данные запроса клиента
type ClientRequest struct {
	Method  string            `json:"method" form:"method" query:"method"`
	URL     string            `json:"url" form:"url" query:"url"`
	Headers map[string]string `json:"headers" form:"headers" query:"headers"`
	Body    string            `json:"body" form:"body" query:"body"`
}

//Do обработка запроса клиента
func (r ClientRequest) Do(timeout time.Duration) (*http.Response, error) {
	netClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	request, err := http.NewRequest(r.Method, r.URL, strings.NewReader(r.Body))
	if err != nil {
		return nil, fmt.Errorf("NewRequest: %v", err)
	}

	for key, value := range r.Headers {
		request.Header.Set(key, value)
	}

	return netClient.Do(request)
}
