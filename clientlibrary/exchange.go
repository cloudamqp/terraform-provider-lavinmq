package clientlibrary

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type ExchangesService struct {
	service
}

type ExchangeRequest struct {
	Type       string         `json:"type,omitempty"`
	AutoDelete *bool          `json:"auto_delete,omitempty"`
	Durable    *bool          `json:"durable,omitempty"`
	Arguments  map[string]any `json:"arguments,omitempty"`
}

type ExchangeResponse struct {
	Name         string                       `json:"name"`
	Vhost        string                       `json:"vhost"`
	Type         string                       `json:"type"`
	AutoDelete   bool                         `json:"auto_delete"`
	Durable      bool                         `json:"durable"`
	Arguments    map[string]any               `json:"arguments,omitempty"`
	MessageStats MessageStatsExchangeResponse `json:"message_stats"`
}

type MessageStatsExchangeResponse struct {
	PublishIn  int64 `json:"publish_in"`
	PublishOut int64 `json:"publish_out"`
	Unroutable int64 `json:"unroutable"`
}

func (s *ExchangesService) CreateOrUpdate(ctx context.Context, vhost, name string, req ExchangeRequest) error {
	path := fmt.Sprintf("api/exchanges/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	tflog.Debug(ctx, s.PathLog("CreateOrUpdate", path))
	_, err := s.client.Request(ctx, http.MethodPut, path, req)
	return err
}

func (s *ExchangesService) Get(ctx context.Context, vhost, name string) (*ExchangeResponse, error) {
	path := fmt.Sprintf("api/exchanges/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	tflog.Debug(ctx, s.PathLog("Get", path))
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result *ExchangeResponse
	err = json.Unmarshal(body, &result)
	tflog.Debug(ctx, s.DataLog("Get", path, result))
	return result, err
}

func (s *ExchangesService) List(ctx context.Context, vhost string) ([]ExchangeResponse, error) {
	path := "api/exchanges"
	if vhost != "" {
		path = fmt.Sprintf("api/exchanges/%s", url.PathEscape(vhost))
	}
	tflog.Debug(ctx, s.PathLog("List", path))
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return []ExchangeResponse{}, err
	}
	if resp == nil {
		return []ExchangeResponse{}, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result []ExchangeResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return []ExchangeResponse{}, err
	}
	return result, nil
}

func (s *ExchangesService) Delete(ctx context.Context, vhost, name string) error {
	path := fmt.Sprintf("api/exchanges/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	tflog.Debug(ctx, s.PathLog("Delete", path))
	_, err := s.client.Request(ctx, http.MethodDelete, path, nil)
	return err
}
