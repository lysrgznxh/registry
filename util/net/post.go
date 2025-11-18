package net

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func PostReq(
	ctx context.Context, httpClient *http.Client,
	apiURL string, request, response interface{}, header map[string]string,
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

	for k, v := range header {
		req.Header.Set(k, v)
	}

	res, err := httpClient.Do(req.WithContext(ctx))
	if res != nil {
		defer (func() { err = res.Body.Close() })()
	}

	if err != nil {
		return err
	}
	//resContent, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(resContent))
	if res.StatusCode != http.StatusOK {
		var errorBody struct {
			Message string `json:"message"`
		}
		if msgerr := json.NewDecoder(res.Body).Decode(&errorBody); msgerr == nil {
			return fmt.Errorf("internal API1: %d from %s: %s", res.StatusCode, apiURL, errorBody.Message)
		}
		resContent, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("response:(%s) status code: %d", resContent, res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(response)
}

func PostForm(
	ctx context.Context, httpClient *http.Client,
	apiURL string, params url.Values, response interface{}, header map[string]string,
) error {

	parsedAPIURL, err := url.Parse(apiURL)
	if err != nil {
		return err
	}

	parsedAPIURL.Path = strings.TrimLeft(parsedAPIURL.Path, "/")
	apiURL = parsedAPIURL.String()

	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer token123")

	for k, v := range header {
		req.Header.Set(k, v)
	}

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

func PostFormV2(
	ctx context.Context, httpClient *http.Client,
	apiURL string, params url.Values, header map[string]string,
) (*http.Response, error) {

	parsedAPIURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}

	parsedAPIURL.Path = strings.TrimLeft(parsedAPIURL.Path, "/")
	apiURL = parsedAPIURL.String()

	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer token123")

	for k, v := range header {
		req.Header.Set(k, v)
	}

	return httpClient.Do(req.WithContext(ctx))
}
