package clientlibrary

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary/utils"
)

type PoliciesService service

type PolicyRequest struct {
	Pattern    string         `json:"pattern"`
	Definition map[string]any `json:"definition"`
	Priority   int64          `json:"priority,omitempty"`
	ApplyTo    string         `json:"apply-to,omitempty"`
}

type PolicyResponse struct {
	Name       string         `json:"name"`
	Vhost      string         `json:"vhost"`
	Pattern    string         `json:"pattern"`
	Definition map[string]any `json:"definition"`
	Priority   int64          `json:"priority"`
	ApplyTo    string         `json:"apply-to"`
}

func (s *PoliciesService) CreateOrUpdate(ctx context.Context, vhost string, name string, policy PolicyRequest) error {
	path := fmt.Sprintf("api/policies/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	_, err := s.client.Request(ctx, http.MethodPut, path, policy)
	return err
}

func (s *PoliciesService) Get(ctx context.Context, vhost string, name string) (*PolicyResponse, error) {
	path := fmt.Sprintf("api/policies/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return utils.GenericUnmarshal[*PolicyResponse](body)
}

func (s *PoliciesService) List(ctx context.Context, vhost string) ([]PolicyResponse, error) {
	path := "api/policies"
	if vhost != "" {
		path = fmt.Sprintf("api/policies/%s", url.PathEscape(vhost))
	}

	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return []PolicyResponse{}, err
	}
	if resp == nil {
		return []PolicyResponse{}, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	result, err := utils.GenericUnmarshal[[]PolicyResponse](body)
	if err != nil {
		return []PolicyResponse{}, err
	}
	return result, nil
}

func (s *PoliciesService) Delete(ctx context.Context, vhost string, name string) error {
	path := fmt.Sprintf("api/policies/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	_, err := s.client.Request(ctx, http.MethodDelete, path, nil)
	return err
}
