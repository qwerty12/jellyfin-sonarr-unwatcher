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
	httpClient *http.Client
	baseUrlUrl *url.URL
	baseUrl    string
	apiKey     string
}

func newSonarrAPIClient(hostUrl string, apiKey string) (*sonarrAPIClient, error) {
	sonarrHostUrl, err := url.Parse(hostUrl)
	if err != nil {
		return nil, err
	}
	if sonarrHostUrl.Scheme == "" || sonarrHostUrl.Host == "" {
		return nil, fmt.Errorf("missing scheme/host")
	}

	sonarrHostUrl = sonarrHostUrl.JoinPath("api", "v3", "/")
	return &sonarrAPIClient{
		baseUrlUrl: sonarrHostUrl,
		baseUrl:    sonarrHostUrl.String(),
		apiKey:     apiKey,
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
	}, nil
}

func (c *sonarrAPIClient) do(method string, endpoint string, queryParams url.Values, reqBody any, respBody any) error {
	var finalUrl string
	if queryParams == nil {
		finalUrl = c.baseUrl + endpoint
	} else {
		u := c.baseUrlUrl.JoinPath(endpoint)
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

func (c *sonarrAPIClient) post(endpoint string, queryParams url.Values, reqBody any, respBody any) error {
	return c.do(http.MethodPost, endpoint, queryParams, reqBody, respBody)
}
