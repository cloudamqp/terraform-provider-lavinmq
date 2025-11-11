package clientlibrary

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type QueuesService service

type QueueRequest struct {
	AutoDelete *bool          `json:"auto_delete,omitempty"`
	Durable    *bool          `json:"durable,omitempty"`
	Arguments  map[string]any `json:"arguments,omitempty"`
}

type QueueResponse struct {
	Name       string         `json:"name"`
	Vhost      string         `json:"vhost"`
	AutoDelete bool           `json:"auto_delete"`
	Durable    bool           `json:"durable"`
	State      string         `json:"state"`
	Consumers  int64          `json:"consumers"`
	Messages   int64          `json:"messages"`
	Ready      int64          `json:"ready"`
	Unacked    int64          `json:"unacked"`
	Arguments  map[string]any `json:"arguments,omitempty"`
}

func (s *QueuesService) CreateOrUpdate(ctx context.Context, vhost, name string, req QueueRequest) error {
	path := fmt.Sprintf("api/queues/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	_, err := s.client.Request(ctx, http.MethodPut, path, req)
	return err
}

func (s *QueuesService) Get(ctx context.Context, vhost, name string) (*QueueResponse, error) {
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
	var result QueueResponse
	err = json.Unmarshal(body, &result)
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
	var result []QueueResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return []QueueResponse{}, err
	}
	return result, nil
}

func (s *QueuesService) Delete(ctx context.Context, vhost, name string) error {
	path := fmt.Sprintf("api/queues/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	_, err := s.client.Request(ctx, http.MethodDelete, path, nil)
	return err
}

func (s *QueuesService) Pause(ctx context.Context, vhost, name string, pause bool) error {
	var path string
	if pause {
		path = fmt.Sprintf("api/queues/%s/%s/pause", url.PathEscape(vhost), url.PathEscape(name))
	} else {
		path = fmt.Sprintf("api/queues/%s/%s/resume", url.PathEscape(vhost), url.PathEscape(name))
	}
	_, err := s.client.Request(ctx, http.MethodPut, path, nil)
	return err
}

func (s *QueuesService) Purge(ctx context.Context, vhost, name string) error {
	path := fmt.Sprintf("api/queues/%s/%s/contents", url.PathEscape(vhost), url.PathEscape(name))
	_, err := s.client.Request(ctx, http.MethodDelete, path, nil)
	return err
}
