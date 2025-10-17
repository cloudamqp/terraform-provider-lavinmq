package clientlibrary

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary/utils"
)

type ShovelsService service

type ShovelResponse struct {
	Name         string  `json:"name"`
	Vhost        string  `json:"vhost"`
	State        string  `json:"state"`
	Error        *string `json:"error,omitempty"`
	MessageCount int64   `json:"message_count"`
}

func (s *ShovelsService) Get(ctx context.Context, vhost, name string) (*ShovelResponse, error) {
	path := fmt.Sprintf("api/shovel/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return utils.GenericUnmarshal[*ShovelResponse](body)
}

func (s *ShovelsService) List(ctx context.Context, vhost string) ([]ShovelResponse, error) {
	path := "api/shovels"
	if vhost != "" {
		path = fmt.Sprintf("api/shovels/%s", url.PathEscape(vhost))
	}
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return []ShovelResponse{}, err
	}
	if resp == nil {
		return []ShovelResponse{}, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	result, err := utils.GenericUnmarshal[[]ShovelResponse](body)
	if err != nil {
		return []ShovelResponse{}, err
	}
	return result, nil
}

func (s *ShovelsService) Pause(ctx context.Context, vhost, name string, pause bool) error {
	var path string
	if pause {
		path = fmt.Sprintf("api/shovels/%s/%s/pause", url.PathEscape(vhost), url.PathEscape(name))
	} else {
		path = fmt.Sprintf("api/shovels/%s/%s/resume", url.PathEscape(vhost), url.PathEscape(name))
	}
	_, err := s.client.Request(ctx, http.MethodPut, path, nil)
	return err
}
