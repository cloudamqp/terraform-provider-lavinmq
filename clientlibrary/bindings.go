package clientlibrary

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type BindingsService service

type BindingRequest struct {
	RoutingKey string         `json:"routing_key"`
	Arguments  map[string]any `json:"arguments,omitempty"`
}

type BindingResponse struct {
	Source          string         `json:"source"`
	Vhost           string         `json:"vhost"`
	Destination     string         `json:"destination"`
	DestinationType string         `json:"destination_type"`
	RoutingKey      string         `json:"routing_key"`
	Arguments       map[string]any `json:"arguments,omitempty"`
	PropertiesKey   string         `json:"properties_key"`
}

func (s *BindingsService) Create(ctx context.Context, vhost, source, destination, destinationType string, req BindingRequest) error {
	var path string
	if destinationType == "queue" {
		path = fmt.Sprintf("api/bindings/%s/e/%s/q/%s", url.PathEscape(vhost), url.PathEscape(source), url.PathEscape(destination))
	} else {
		path = fmt.Sprintf("api/bindings/%s/e/%s/e/%s", url.PathEscape(vhost), url.PathEscape(source), url.PathEscape(destination))
	}
	_, err := s.client.Request(ctx, http.MethodPost, path, req)
	return err
}

func (s *BindingsService) Get(ctx context.Context, vhost, source, destination, destinationType, propertiesKey string) (*BindingResponse, error) {
	var path string
	if destinationType == "queue" {
		path = fmt.Sprintf("api/bindings/%s/e/%s/q/%s/%s", url.PathEscape(vhost), url.PathEscape(source), url.PathEscape(destination), url.PathEscape(propertiesKey))
	} else {
		path = fmt.Sprintf("api/bindings/%s/e/%s/e/%s/%s", url.PathEscape(vhost), url.PathEscape(source), url.PathEscape(destination), url.PathEscape(propertiesKey))
	}
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result *BindingResponse
	err = json.Unmarshal(body, &result)
	return result, err
}

func (s *BindingsService) List(ctx context.Context, vhost string) ([]BindingResponse, error) {
	path := "api/bindings"
	if vhost != "" {
		path = fmt.Sprintf("api/bindings/%s", url.PathEscape(vhost))
	}
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return []BindingResponse{}, err
	}
	if resp == nil {
		return []BindingResponse{}, nil
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result []BindingResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return []BindingResponse{}, err
	}
	return result, nil
}

func (s *BindingsService) Delete(ctx context.Context, vhost, source, destination, destinationType, propertiesKey string) error {
	var path string
	if destinationType == "queue" {
		path = fmt.Sprintf("api/bindings/%s/e/%s/q/%s/%s", url.PathEscape(vhost), url.PathEscape(source), url.PathEscape(destination), url.PathEscape(propertiesKey))
	} else {
		path = fmt.Sprintf("api/bindings/%s/e/%s/e/%s/%s", url.PathEscape(vhost), url.PathEscape(source), url.PathEscape(destination), url.PathEscape(propertiesKey))
	}
	_, err := s.client.Request(ctx, http.MethodDelete, path, nil)
	return err
}
