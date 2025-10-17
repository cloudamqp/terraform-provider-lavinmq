package clientlibrary

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary/utils"
)

type ShovelParametersService service

type ShovelParametersRequest struct {
	Value ShovelParametersObject `json:"value"`
}

type ShovelParametersResponse struct {
	Name  string                 `json:"name"`
	Vhost string                 `json:"vhost"`
	Value ShovelParametersObject `json:"value"`
}

type ShovelParametersObject struct {
	SrcUri  string `json:"src-uri"`
	DestUri string `json:"dest-uri"`
	// Optional parameters
	SrcQueue         *string `json:"src-queue,omitempty"`
	SrcExchange      *string `json:"src-exchange,omitempty"`
	SrcExchangeKey   *string `json:"src-exchange-key,omitempty"`
	SrcPrefetchCount *int64  `json:"src-prefetch-count,omitempty"`
	SrcDelayAfter    *string `json:"src-delay-after,omitempty"`
	DestQueue        *string `json:"dest-queue,omitempty"`
	DestExchange     *string `json:"dest-exchange,omitempty"`
	DestExchangeKey  *string `json:"dest-exchange-key,omitempty"`
	ReconnectDelay   *int64  `json:"reconnect-delay,omitempty"`
	AckMode          *string `json:"ack-mode,omitempty"`
}

func (s *ShovelParametersService) CreateOrUpdate(ctx context.Context, vhost, name string, shovel ShovelParametersRequest) error {
	path := fmt.Sprintf("api/parameters/shovel/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	_, err := s.client.Request(ctx, http.MethodPut, path, shovel)
	return err
}

func (s *ShovelParametersService) Get(ctx context.Context, vhost, name string) (*ShovelParametersResponse, error) {
	path := fmt.Sprintf("api/parameters/shovel/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return utils.GenericUnmarshal[*ShovelParametersResponse](body)
}

func (s *ShovelParametersService) Delete(ctx context.Context, vhost, name string) error {
	path := fmt.Sprintf("api/parameters/shovel/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	_, err := s.client.Request(ctx, http.MethodDelete, path, nil)
	return err
}

func (s *ShovelParametersService) List(ctx context.Context, vhost string) ([]ShovelParametersResponse, error) {
	path := "api/parameters/shovel"
	if vhost != "" {
		path = fmt.Sprintf("api/parameters/shovel/%s", url.PathEscape(vhost))
	}
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return []ShovelParametersResponse{}, err
	}
	if resp == nil {
		return []ShovelParametersResponse{}, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	result, err := utils.GenericUnmarshal[[]ShovelParametersResponse](body)
	if err != nil {
		return []ShovelParametersResponse{}, err
	}
	return result, nil
}
