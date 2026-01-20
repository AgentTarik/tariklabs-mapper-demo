package engine

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HTTPEngine interface {
	MakeRequest(baseURL string, method string, body []byte, params map[string]string) ([]byte, int, error)
}

type httpEngine struct {
	client *http.Client
}

func NewHTTPEngine() HTTPEngine {
	return &httpEngine{
		client: &http.Client{},
	}
}

func (e *httpEngine) MakeRequest(baseURL string, method string, body []byte, params map[string]string) ([]byte, int, error) {
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid url: %w", err)
	}

	if params != nil {
		query := reqURL.Query()
		for key, value := range params {
			query.Set(key, value)
		}
		reqURL.RawQuery = query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(method, reqURL.String(), bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response: %w", err)
	}

	return respBody, resp.StatusCode, nil
}
