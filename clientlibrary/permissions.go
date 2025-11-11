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

type PermissionsService service

type PermissionRequest struct {
	Configure string `json:"configure"`
	Read      string `json:"read"`
	Write     string `json:"write"`
}

type PermissionResponse struct {
	User      string `json:"user"`
	Vhost     string `json:"vhost"`
	Configure string `json:"configure"`
	Read      string `json:"read"`
	Write     string `json:"write"`
}

func (s *PermissionsService) CreateOrUpdate(ctx context.Context, vhost, user string, permission PermissionRequest) error {
	path := fmt.Sprintf("api/permissions/%s/%s", url.PathEscape(vhost), url.PathEscape(user))
	tflog.Debug(ctx, fmt.Sprintf("service=permissions method=CreateOrUpdate path=%s", path))
	_, err := s.client.Request(ctx, http.MethodPut, path, permission)
	return err
}

func (s *PermissionsService) Get(ctx context.Context, vhost, user string) (*PermissionResponse, error) {
	path := fmt.Sprintf("api/permissions/%s/%s", url.PathEscape(vhost), url.PathEscape(user))
	tflog.Debug(ctx, fmt.Sprintf("service=permissions method=Get path=%s", path))
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result *PermissionResponse
	err = json.Unmarshal(body, &result)
	tflog.Debug(ctx, fmt.Sprintf("service=permissions method=Get path=%s, result=%+v", path, result))
	return result, err
}

func (s *PermissionsService) List(ctx context.Context, vhost, user string) ([]PermissionResponse, error) {
	var path string
	switch {
	case vhost != "" && user != "":
		path = fmt.Sprintf("api/vhosts/%s/permissions", url.PathEscape(vhost))
	case vhost != "":
		path = fmt.Sprintf("api/vhosts/%s/permissions", url.PathEscape(vhost))
	case user != "":
		path = fmt.Sprintf("api/users/%s/permissions", url.PathEscape(user))
	default:
		path = "api/permissions"
	}
	tflog.Debug(ctx, fmt.Sprintf("service=permissions method=List path=%s", path))

	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return []PermissionResponse{}, err
	}
	if resp == nil {
		return []PermissionResponse{}, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result []PermissionResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return []PermissionResponse{}, err
	}
	return result, nil
}

func (s *PermissionsService) Delete(ctx context.Context, vhost, user string) error {
	path := fmt.Sprintf("api/permissions/%s/%s", url.PathEscape(vhost), url.PathEscape(user))
	tflog.Debug(ctx, fmt.Sprintf("service=permissions method=Delete path=%s", path))
	_, err := s.client.Request(ctx, http.MethodDelete, path, nil)
	return err
}
