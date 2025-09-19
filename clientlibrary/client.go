package clientlibrary

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary/utils"
)

type Client struct {
	HttpClient *http.Client
	BaseURL    string
	UserAgent  string

	Username string
	Password string

	common      service
	Users       *UsersService
	VhostLimits *VhostLimitsService
	Vhosts      *VhostsService
}

type service struct {
	client *Client
}

type ErrorResponse struct {
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

func NewClient(baseURL, useragent, username, password string, httpClient *http.Client) *Client {
	client := &Client{
		BaseURL:    baseURL,
		UserAgent:  useragent,
		HttpClient: httpClient,
		Username:   username,
		Password:   password,
	}
	client.initialize()
	return client
}

func (c *Client) initialize() {
	c.common.client = c
	c.Users = (*UsersService)(&c.common)
	c.VhostLimits = (*VhostLimitsService)(&c.common)
	c.Vhosts = (*VhostsService)(&c.common)
}

func (c *Client) NewRequest(method, path string, body any) (*http.Request, error) {
	url := c.BaseURL + path

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	resp, err := c.HttpClient.Do(req)
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
	default:
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		errorBody, _ := utils.GenericUnmarshal[ErrorResponse](body)
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
