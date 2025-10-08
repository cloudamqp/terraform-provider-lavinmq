package clientlibrary

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary/utils"
)

type QueuesService service

type QueueRequest struct {
	AutoDelete *bool                  `json:"auto_delete,omitempty"`
	Durable    *bool                  `json:"durable,omitempty"`
	Arguments  map[string]interface{} `json:"arguments,omitempty"`
}

type QueueResponse struct {
	Name       string                 `json:"name"`
	Vhost      string                 `json:"vhost"`
	AutoDelete bool                   `json:"auto_delete"`
	Durable    bool                   `json:"durable"`
	Arguments  map[string]interface{} `json:"arguments,omitempty"`
}

func (s *QueuesService) CreateOrUpdate(ctx context.Context, vhost string, name string, req QueueRequest) error {
	path := fmt.Sprintf("api/queues/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	_, err := s.client.Request(ctx, http.MethodPut, path, req)
	return err
}

func (s *QueuesService) Get(ctx context.Context, vhost string, name string) (*QueueResponse, error) {
	path := fmt.Sprintf("api/queues/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	result, err := utils.GenericUnmarshal[QueueResponse](body)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *QueuesService) List(ctx context.Context, vhost string) ([]QueueResponse, error) {
	path := "api/queues"
	if vhost != "" {
		path = fmt.Sprintf("api/queues/%s", url.PathEscape(vhost))
	}
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return []QueueResponse{}, err
	}
	if resp == nil {
		return []QueueResponse{}, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	result, err := utils.GenericUnmarshal[[]QueueResponse](body)
	if err != nil {
		return []QueueResponse{}, err
	}
	return result, nil
}

func (s *QueuesService) Delete(ctx context.Context, vhost string, name string) error {
	path := fmt.Sprintf("api/queues/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	_, err := s.client.Request(ctx, http.MethodDelete, path, nil)
	return err
}
