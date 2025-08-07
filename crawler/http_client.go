package crawler

import (
	"fmt"
	"net/http"
	"time"
)

func NewHTTPClient(timeoutSec int) http.Client {
	return http.Client{
		Timeout: time.Duration(timeoutSec) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		},
	}
}

func NewRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64)")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	return req, nil
}
