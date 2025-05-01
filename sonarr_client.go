package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type sonarrAPIClient struct {
	url        string
	apiKey     string
	httpClient *http.Client
}

func newSonarrAPIClient(baseUrl *url.URL, apiKey string) *sonarrAPIClient {
	return &sonarrAPIClient{
		url:    baseUrl.JoinPath("api", "v3", "/").String(),
		apiKey: apiKey,
		httpClient: &http.Client{
			Transport: &http.Transport{
				Proxy:                 nil, // $HTTP_PROXY etc. ignored
				MaxIdleConns:          50,
				IdleConnTimeout:       time.Minute,
				TLSHandshakeTimeout:   http.DefaultTransport.(*http.Transport).TLSHandshakeTimeout,
				ExpectContinueTimeout: http.DefaultTransport.(*http.Transport).ExpectContinueTimeout,
				ResponseHeaderTimeout: 10 * time.Second,
				DialContext:           (&net.Dialer{Timeout: 30 * time.Second, KeepAlive: time.Minute}).DialContext,
				ForceAttemptHTTP2:     false,
			},
		},
	}
}

func (c *sonarrAPIClient) do(method string, endpoint string, queryParams url.Values, reqBody any, respBody any) error {
	finalUrl := c.url + endpoint
	if queryParams != nil {
		u, err := url.Parse(finalUrl)
		if err != nil {
			return err
		}

		u.RawQuery = queryParams.Encode()

		finalUrl = u.String()
	}

	var pReqBody io.Reader = nil
	var jsonBuf bytes.Buffer
	if reqBody != nil {
		jsonEnc := json.NewEncoder(&jsonBuf)
		jsonEnc.SetEscapeHTML(false)
		if err := jsonEnc.Encode(reqBody); err != nil {
			return fmt.Errorf("failed to serialise request body to JSON for %s: %w", finalUrl, err)
		}
		pReqBody = &jsonBuf
	}

	req, err := http.NewRequest(method, finalUrl, pReqBody)
	if err != nil {
		return fmt.Errorf("failed to create %s request for %s: %w", method, finalUrl, err)
	}
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if respBody != nil {
		req.Header.Set("Accept", "application/json")
	}
	req.Header.Set("X-Api-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute %s request for %s: %w", method, finalUrl, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to %s %s: %s", method, finalUrl, resp.Status)
	}

	if respBody != nil {
		if err = json.NewDecoder(resp.Body).Decode(respBody); err != nil {
			return fmt.Errorf("failed to decode JSON response from %s: %w", finalUrl, err)
		}
	}

	return nil
}

func (c *sonarrAPIClient) get(endpoint string, queryParams url.Values, respBody any) error {
	return c.do(http.MethodGet, endpoint, queryParams, nil, respBody)
}

func (c *sonarrAPIClient) put(endpoint string, queryParams url.Values, reqBody any, respBody any) error {
	return c.do(http.MethodPut, endpoint, queryParams, reqBody, respBody)
}
