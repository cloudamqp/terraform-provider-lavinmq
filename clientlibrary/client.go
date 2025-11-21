package clientlibrary

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string

	username string
	password string
}

type service struct {
	client *Client
	name   string
}

func (s *service) PathLog(method, path string) string {
	return fmt.Sprintf("service=%s method=%s path=%s", s.name, method, path)
}

func (s *service) DataLog(method, path string, data any) string {
	return fmt.Sprintf("service=%s method=%s path=%s data=%+v", s.name, method, path, data)
}

type ErrorResponse struct {
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

func NewClient(baseURL, useragent, username, password string, httpClient *http.Client) *Client {
	return &Client{
		baseURL:    baseURL,
		userAgent:  useragent,
		httpClient: httpClient,
		username:   username,
		password:   password,
	}
}

func (c *Client) NewRequest(method, path string, body any) (*http.Request, error) {
	baseURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	// Parse the path and join it with the base URL
	fullURL, err := url.JoinPath(baseURL.String(), path)
	if err != nil {
		return nil, fmt.Errorf("failed to join URL path: %w", err)
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, fullURL, buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}

	switch resp.StatusCode {
	case 200, 201, 204:
		return resp, nil
	case 404:
		return nil, nil
	default:
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var errorBody ErrorResponse
		_ = json.Unmarshal(body, &errorBody)
		return nil, fmt.Errorf("status code: %d, error: %s", resp.StatusCode, errorBody.Reason)
	}
}

func (c *Client) Request(ctx context.Context, method, path string, body any) (*http.Response, error) {
	req, err := c.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req)
}
