package clientlibrary

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary/utils"
)

type UsersService service

type UserRequest struct {
	Password     *string `json:"password,omitempty"`
	PasswordHash *string `json:"password_hash,omitempty"`
	Tags         string  `json:"tags"`
}

type UserResponse struct {
	Name         string  `json:"name"`
	PasswordHash *string `json:"password_hash,omitempty"`
	Tags         string  `json:"tags"`
}

func (s *UsersService) CreateOrUpdate(ctx context.Context, username string, user UserRequest) error {
	path := fmt.Sprintf("api/users/%s", url.PathEscape(username))
	_, err := s.client.Request(ctx, http.MethodPut, path, user)
	return err
}

func (s *UsersService) Get(ctx context.Context, username string) (*UserResponse, error) {
	path := fmt.Sprintf("api/users/%s", url.PathEscape(username))
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return utils.GenericUnmarshal[*UserResponse](body)
}

func (s *UsersService) List(ctx context.Context) ([]UserResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "api/users", nil)
	if err != nil {
		return []UserResponse{}, err
	}
	if resp == nil {
		return []UserResponse{}, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	result, err := utils.GenericUnmarshal[[]UserResponse](body)
	if err != nil {
		return []UserResponse{}, err
	}
	return result, nil
}

func (s *UsersService) Delete(ctx context.Context, username string) error {
	path := fmt.Sprintf("api/users/%s", url.PathEscape(username))
	_, err := s.client.Request(ctx, http.MethodDelete, path, nil)
	return err
}
