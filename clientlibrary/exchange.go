package clientlibrary

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary/utils"
)

type ExchangesService service

type ExchangeRequest struct {
	Type       string                 `json:"type,omitempty"`
	AutoDelete *bool                  `json:"auto_delete,omitempty"`
	Durable    *bool                  `json:"durable,omitempty"`
	Internal   *bool                  `json:"internal,omitempty"`
	Arguments  map[string]interface{} `json:"arguments,omitempty"`
}

type ExchangeResponse struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	AutoDelete bool                   `json:"auto_delete"`
	Durable    bool                   `json:"durable"`
	Internal   bool                   `json:"internal"`
	Arguments  map[string]interface{} `json:"arguments,omitempty"`
}

func (s *ExchangesService) CreateOrUpdate(ctx context.Context, vhost string, name string, req ExchangeRequest) error {
	path := fmt.Sprintf("api/exchanges/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	_, err := s.client.Request(ctx, http.MethodPut, path, req)
	return err
}

func (s *ExchangesService) Get(ctx context.Context, vhost string, name string) (*ExchangeResponse, error) {
	path := fmt.Sprintf("api/exchanges/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return utils.GenericUnmarshal[*ExchangeResponse](body)
}

func (s *ExchangesService) Delete(ctx context.Context, vhost string, name string) error {
	path := fmt.Sprintf("api/exchanges/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	_, err := s.client.Request(ctx, http.MethodDelete, path, nil)
	return err
}
