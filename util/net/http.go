package net

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// PostJSON performs a POST request with JSON on an internal HTTP API
func Post(
	ctx context.Context, httpClient *http.Client,
	apiURL string, request interface{},
) (*http.Response, error) {
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	parsedAPIURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}

	parsedAPIURL.Path = strings.TrimLeft(parsedAPIURL.Path, "/")
	apiURL = parsedAPIURL.String()

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return httpClient.Do(req.WithContext(ctx))
}

// PostJSON performs a POST request with JSON on an internal HTTP API
func PostJSON(
	ctx context.Context, httpClient *http.Client,
	apiURL string, request, response interface{},
) error {
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	parsedAPIURL, err := url.Parse(apiURL)
	if err != nil {
		return err
	}

	parsedAPIURL.Path = strings.TrimLeft(parsedAPIURL.Path, "/")
	apiURL = parsedAPIURL.String()

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := httpClient.Do(req.WithContext(ctx))
	if res != nil {
		defer (func() { err = res.Body.Close() })()
	}

	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		var errorBody struct {
			Message string `json:"message"`
		}
		if msgerr := json.NewDecoder(res.Body).Decode(&errorBody); msgerr == nil {
			return fmt.Errorf("internal API: %d from %s: %s", res.StatusCode, apiURL, errorBody.Message)
		}
		return fmt.Errorf("internal API: %d from %s", res.StatusCode, apiURL)
	}
	return json.NewDecoder(res.Body).Decode(response)
}
