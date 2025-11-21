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

type UsersService struct {
	service
}

type UserRequest struct {
	Password         string `json:"password,omitempty"`
	PasswordHash     string `json:"password_hash,omitempty"`
	HashingAlgorithm string `json:"hashing_algorithm,omitempty"`
	Tags             string `json:"tags"`
}

type UserResponse struct {
	Name             string `json:"name"`
	PasswordHash     string `json:"password_hash"`
	HashingAlgorithm string `json:"hashing_algorithm"`
	Tags             string `json:"tags"`
}

func (u UserRequest) Sanitized() UserRequest {
	sanitized := u
	if sanitized.Password != "" {
		sanitized.Password = "***"
	}
	if sanitized.PasswordHash != "" {
		sanitized.PasswordHash = "***"
	}
	return sanitized
}

func (u UserResponse) Sanitized() UserResponse {
	sanitized := u
	if sanitized.PasswordHash != "" {
		sanitized.PasswordHash = "***"
	}
	return sanitized
}

func (s *UsersService) CreateOrUpdate(ctx context.Context, username string, user UserRequest) error {
	path := fmt.Sprintf("api/users/%s", url.PathEscape(username))
	tflog.Debug(ctx, s.DataLog("CreateOrUpdate", path, user.Sanitized()))
	_, err := s.client.Request(ctx, http.MethodPut, path, user)
	return err
}

func (s *UsersService) Get(ctx context.Context, username string) (*UserResponse, error) {
	path := fmt.Sprintf("api/users/%s", url.PathEscape(username))
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
	var result *UserResponse
	err = json.Unmarshal(body, &result)
	tflog.Debug(ctx, s.DataLog("Get", path, result.Sanitized()))
	return result, err
}

func (s *UsersService) List(ctx context.Context) ([]UserResponse, error) {
	tflog.Debug(ctx, s.PathLog("List", "api/users"))
	resp, err := s.client.Request(ctx, http.MethodGet, "api/users", nil)
	if err != nil {
		return []UserResponse{}, err
	}
	if resp == nil {
		return []UserResponse{}, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result []UserResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return []UserResponse{}, err
	}
	return result, nil
}

func (s *UsersService) Delete(ctx context.Context, username string) error {
	path := fmt.Sprintf("api/users/%s", url.PathEscape(username))
	tflog.Debug(ctx, s.PathLog("Delete", path))
	_, err := s.client.Request(ctx, http.MethodDelete, path, nil)
	return err
}
