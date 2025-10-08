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
	Arguments  map[string]interface{} `json:"arguments,omitempty"`
}

type ExchangeResponse struct {
	Name       string                 `json:"name"`
	Vhost      string                 `json:"vhost"`
	Type       string                 `json:"type"`
	AutoDelete bool                   `json:"auto_delete"`
	Durable    bool                   `json:"durable"`
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

func (s *ExchangesService) List(ctx context.Context, vhost string) ([]*ExchangeResponse, error) {
	path := "api/exchanges"
	if vhost != "" {
		path = fmt.Sprintf("api/exchanges/%s", url.PathEscape(vhost))
	}
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	result, err := utils.GenericUnmarshal[[]ExchangeResponse](body)
	if err != nil {
		return nil, err
	}

	// Convert slice of values to slice of pointers
	pointers := make([]*ExchangeResponse, len(result))
	for i := range result {
		pointers[i] = &result[i]
	}
	return pointers, nil
}

func (s *ExchangesService) Delete(ctx context.Context, vhost string, name string) error {
	path := fmt.Sprintf("api/exchanges/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	_, err := s.client.Request(ctx, http.MethodDelete, path, nil)
	return err
}
