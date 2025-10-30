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

type VhostLimitsService service

type vhostLimitValueRequest struct {
	Value int64 `json:"value"`
}

type VhostLimitsResponse struct {
	Vhost string      `json:"vhost"`
	Value VhostLimits `json:"value"`
}

type VhostLimits struct {
	MaxConnections *int64 `json:"max-connections,omitempty"`
	MaxQueues      *int64 `json:"max-queues,omitempty"`
}

func (s *VhostLimitsService) Update(ctx context.Context, vhost string, limits VhostLimits) error {
	if limits.MaxConnections != nil {
		path := fmt.Sprintf("/api/vhost-limits/%s/%s", url.PathEscape(vhost), "max-connections")
		value := vhostLimitValueRequest{Value: *limits.MaxConnections}
		tflog.Debug(ctx, fmt.Sprintf("Vhost limits update path: %s, value: %v\n", path, value))
		if _, err := s.client.Request(ctx, http.MethodPut, path, value); err != nil {
			return err
		}
	} else {
		_ = s.Delete(ctx, vhost, "max-connections")
	}

	if limits.MaxQueues != nil {
		path := fmt.Sprintf("/api/vhost-limits/%s/%s", url.PathEscape(vhost), "max-queues")
		value := vhostLimitValueRequest{Value: *limits.MaxQueues}
		tflog.Debug(ctx, fmt.Sprintf("Vhost limits update path: %s, value: %v\n", path, value))
		if _, err := s.client.Request(ctx, http.MethodPut, path, value); err != nil {
			return err
		}
	} else {
		_ = s.Delete(ctx, vhost, "max-queues")
	}

	return nil
}

func (s *VhostLimitsService) Get(ctx context.Context, vhost string) (VhostLimitsResponse, error) {
	path := fmt.Sprintf("api/vhost-limits/%s", url.PathEscape(vhost))
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return VhostLimitsResponse{}, err
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var vhostLimitsResponses []VhostLimitsResponse
	err = json.Unmarshal(body, &vhostLimitsResponses)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Unmarshal failed: %v", err))
		return VhostLimitsResponse{}, err
	}
	for _, v := range vhostLimitsResponses {
		if v.Vhost == vhost {
			return v, nil
		}
	}
	return VhostLimitsResponse{}, nil
}

func (s *VhostLimitsService) Delete(ctx context.Context, vhost, limitType string) error {
	path := fmt.Sprintf("api/vhost-limits/%s/%s", url.PathEscape(vhost), url.PathEscape(limitType))
	tflog.Debug(ctx, fmt.Sprintf("Remove limit type: %s for vhost: %s", limitType, vhost))
	_, err := s.client.Request(ctx, http.MethodDelete, path, nil)
	return err
}
