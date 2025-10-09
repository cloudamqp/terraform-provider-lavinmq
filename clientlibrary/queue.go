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
	Arguments  map[string]any `json:"arguments,omitempty"`
}

type QueueResponse struct {
	AutoDelete bool                   `json:"auto_delete"`
	Durable    bool                   `json:"durable"`
	Arguments  map[string]any `json:"arguments,omitempty"`
}

func (s *QueuesService) CreateOrUpdate(ctx context.Context, vhost string, name string, req QueueRequest) error {
	path := fmt.Sprintf("api/queues/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	_, err := s.client.Request(ctx, http.MethodPut, path, req)
	return err
}

func (s *QueuesService) Get(ctx context.Context, vhost string, name string) (QueueResponse, error) {
	path := fmt.Sprintf("api/queues/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return QueueResponse{}, err
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return utils.GenericUnmarshal[QueueResponse](body)
}

func (s *QueuesService) Delete(ctx context.Context, vhost string, name string) error {
	path := fmt.Sprintf("api/queues/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	_, err := s.client.Request(ctx, http.MethodDelete, path, nil)
	return err
}
