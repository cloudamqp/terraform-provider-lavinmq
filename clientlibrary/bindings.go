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
	vhost = url.PathEscape(vhost)
	source = url.PathEscape(source)
	destination = url.PathEscape(destination)

	var path string
	if destinationType == "queue" {
		path = fmt.Sprintf("api/bindings/%s/e/%s/q/%s", vhost, source, destination)
	} else {
		path = fmt.Sprintf("api/bindings/%s/e/%s/e/%s", vhost, source, destination)
	}
	tflog.Debug(ctx, fmt.Sprintf("service=bindings method=Create path=%s", path))
	_, err := s.client.Request(ctx, http.MethodPost, path, req)
	return err
}

func (s *BindingsService) Get(ctx context.Context, vhost, source, destination, destinationType, propertiesKey string) (*BindingResponse, error) {
	vhost = url.PathEscape(vhost)
	source = url.PathEscape(source)
	destination = url.PathEscape(destination)
	propertiesKey = url.PathEscape(propertiesKey)

	var path string
	if destinationType == "queue" {
		path = fmt.Sprintf("api/bindings/%s/e/%s/q/%s/%s", vhost, source, destination, propertiesKey)
	} else {
		path = fmt.Sprintf("api/bindings/%s/e/%s/e/%s/%s", vhost, source, destination, propertiesKey)
	}
	tflog.Debug(ctx, fmt.Sprintf("service=bindings method=Get path=%s", path))
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
	tflog.Debug(ctx, fmt.Sprintf("service=bindings method=Get path=%s, result=%+v", path, result))
	return result, err
}

func (s *BindingsService) List(ctx context.Context, vhost string) ([]BindingResponse, error) {
	path := "api/bindings"
	if vhost != "" {
		path = fmt.Sprintf("api/bindings/%s", url.PathEscape(vhost))
	}
	tflog.Debug(ctx, fmt.Sprintf("service=bindings method=List path=%s", path))
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
	vhost = url.PathEscape(vhost)
	source = url.PathEscape(source)
	destination = url.PathEscape(destination)
	propertiesKey = url.PathEscape(propertiesKey)

	var path string
	if destinationType == "queue" {
		path = fmt.Sprintf("api/bindings/%s/e/%s/q/%s/%s", vhost, source, destination, propertiesKey)
	} else {
		path = fmt.Sprintf("api/bindings/%s/e/%s/e/%s/%s", vhost, source, destination, propertiesKey)
	}
	tflog.Debug(ctx, fmt.Sprintf("service=bindings method=Delete path=%s", path))
	_, err := s.client.Request(ctx, http.MethodDelete, path, nil)
	return err
}
