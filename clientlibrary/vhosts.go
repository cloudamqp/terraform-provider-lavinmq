package clientlibrary

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type VhostsService service

type VhostResponse struct {
	Name                   string                    `json:"name"`
	Dir                    string                    `json:"dir"`
	Tracing                bool                      `json:"tracing"`
	Messages               int64                     `json:"messages"`
	MessagesUnacknowledged int64                     `json:"messages_unacknowledged"`
	MessagesReady          int64                     `json:"messages_ready"`
	MessagesStats          VhostMessageStatsResponse `json:"messages_stats"`
}

type VhostMessageStatsResponse struct {
	Ack              int64 `json:"ack"`
	Confirm          int64 `json:"confirm"`
	Deliver          int64 `json:"deliver"`
	Get              int64 `json:"get"`
	GetNoAck         int64 `json:"get_no_ack"`
	Publish          int64 `json:"publish"`
	Redeliver        int64 `json:"redeliver"`
	ReturnUnroutable int64 `json:"return_unroutable"`
}

func (s *VhostsService) CreateOrUpdate(ctx context.Context, name string) error {
	path := fmt.Sprintf("api/vhosts/%s", url.PathEscape(name))
	_, err := s.client.Request(ctx, http.MethodPut, path, nil)
	return err
}

func (s *VhostsService) Get(ctx context.Context, name string) (VhostResponse, error) {
	path := fmt.Sprintf("api/vhosts/%s", url.PathEscape(name))
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return VhostResponse{}, err
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result VhostResponse
	err = json.Unmarshal(body, &result)
	return result, err
}

func (s *VhostsService) List(ctx context.Context) ([]VhostResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "api/vhosts", nil)
	if err != nil {
		return []VhostResponse{}, err
	}
	if resp == nil {
		return []VhostResponse{}, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result []VhostResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return []VhostResponse{}, err
	}
	return result, nil
}

func (s *VhostsService) Delete(ctx context.Context, name string) error {
	path := fmt.Sprintf("api/vhosts/%s", url.PathEscape(name))
	_, err := s.client.Request(ctx, http.MethodDelete, path, nil)
	return err
}
