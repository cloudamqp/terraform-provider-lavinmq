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

type PoliciesService struct {
	service
}

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

func (s *PoliciesService) CreateOrUpdate(ctx context.Context, vhost, name string, policy PolicyRequest) error {
	path := fmt.Sprintf("api/policies/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	tflog.Debug(ctx, s.PathLog("CreateOrUpdate", path))
	_, err := s.client.Request(ctx, http.MethodPut, path, policy)
	return err
}

func (s *PoliciesService) Get(ctx context.Context, vhost, name string) (*PolicyResponse, error) {
	path := fmt.Sprintf("api/policies/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
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
	var result *PolicyResponse
	err = json.Unmarshal(body, &result)
	tflog.Debug(ctx, s.DataLog("Get", path, result))
	return result, err
}

func (s *PoliciesService) List(ctx context.Context, vhost string) ([]PolicyResponse, error) {
	path := "api/policies"
	if vhost != "" {
		path = fmt.Sprintf("api/policies/%s", url.PathEscape(vhost))
	}
	tflog.Debug(ctx, s.PathLog("List", path))

	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return []PolicyResponse{}, err
	}
	if resp == nil {
		return []PolicyResponse{}, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result []PolicyResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return []PolicyResponse{}, err
	}
	return result, nil
}

func (s *PoliciesService) Delete(ctx context.Context, vhost, name string) error {
	path := fmt.Sprintf("api/policies/%s/%s", url.PathEscape(vhost), url.PathEscape(name))
	tflog.Debug(ctx, s.PathLog("Delete", path))
	_, err := s.client.Request(ctx, http.MethodDelete, path, nil)
	return err
}
