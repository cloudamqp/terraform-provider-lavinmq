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

type ParametersService service

type ShovelValue struct {
	SrcURI           string `json:"src-uri"`
	DestURI          string `json:"dest-uri"`
	SrcQueue         string `json:"src-queue,omitempty"`
	SrcExchange      string `json:"src-exchange,omitempty"`
	SrcExchangeKey   string `json:"src-exchange-key,omitempty"`
	DestQueue        string `json:"dest-queue,omitempty"`
	DestExchange     string `json:"dest-exchange,omitempty"`
	DestExchangeKey  string `json:"dest-exchange-key,omitempty"`
	SrcPrefetchCount int64  `json:"src-prefetch-count,omitempty"`
	SrcDeleteAfter   string `json:"src-delete-after,omitempty"`
	ReconnectDelay   int64  `json:"reconnect-delay,omitempty"`
	AckMode          string `json:"ack-mode,omitempty"`
}

type FederationUpstreamValue struct {
	URI            string `json:"uri"`
	PrefetchCount  int64  `json:"prefetch-count,omitempty"`
	ReconnectDelay int64  `json:"reconnect-delay,omitempty"`
	AckMode        string `json:"ack-mode,omitempty"`
	Exchange       string `json:"exchange,omitempty"`
	MaxHops        int64  `json:"max-hops,omitempty"`
	Expires        int64  `json:"expires,omitempty"`
	MessageTTL     int64  `json:"message-ttl,omitempty"`
	Queue          string `json:"queue,omitempty"`
	ConsumerTag    string `json:"consumer-tag,omitempty"`
}

type ParameterRequest struct {
	Value any `json:"value"`
}

type ParameterResponse struct {
	Name      string `json:"name"`
	Vhost     string `json:"vhost"`
	Component string `json:"component"`
	Value     any    `json:"value"`
}

func (s *ParametersService) CreateOrUpdate(ctx context.Context, component, vhost, name string, request ParameterRequest) error {
	path := fmt.Sprintf("api/parameters/%s/%s/%s", url.PathEscape(component), url.PathEscape(vhost), url.PathEscape(name))
	tflog.Debug(ctx, fmt.Sprintf("service=parameters method=CreateOrUpdate path=%s", path))
	_, err := s.client.Request(ctx, http.MethodPut, path, request)
	return err
}

func (s *ParametersService) Get(ctx context.Context, component, vhost, name string) (*ParameterResponse, error) {
	path := fmt.Sprintf("api/parameters/%s/%s/%s", url.PathEscape(component), url.PathEscape(vhost), url.PathEscape(name))
	tflog.Debug(ctx, fmt.Sprintf("service=parameters method=Get path=%s", path))
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result *ParameterResponse
	err = json.Unmarshal(body, &result)
	tflog.Debug(ctx, fmt.Sprintf("service=parameters method=Get path=%s, result=%+v", path, result))
	return result, err
}

func (s *ParametersService) List(ctx context.Context, component, vhost string) ([]ParameterResponse, error) {
	path := "api/parameters"
	if component != "" {
		component = url.PathEscape(component)
		path = fmt.Sprintf("api/parameters/%s", component)
		if vhost != "" {
			vhost = url.PathEscape(vhost)
			path = fmt.Sprintf("api/parameters/%s/%s", component, vhost)
		}
	}
	tflog.Debug(ctx, fmt.Sprintf("service=parameters method=List path=%s", path))

	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return []ParameterResponse{}, err
	}
	if resp == nil {
		return []ParameterResponse{}, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result []ParameterResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return []ParameterResponse{}, err
	}
	return result, nil
}

func (s *ParametersService) Delete(ctx context.Context, component, vhost, name string) error {
	path := fmt.Sprintf("api/parameters/%s/%s/%s", url.PathEscape(component), url.PathEscape(vhost), url.PathEscape(name))
	tflog.Debug(ctx, fmt.Sprintf("service=parameters method=Delete path=%s", path))
	_, err := s.client.Request(ctx, http.MethodDelete, path, nil)
	return err
}
