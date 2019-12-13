package request

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
)

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
