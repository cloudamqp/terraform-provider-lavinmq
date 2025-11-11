package clientlibrary

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type MessagesService service

type PublishRequest struct {
	RoutingKey      string         `json:"routing_key"`
	Payload         string         `json:"payload"`
	PayloadEncoding string         `json:"payload_encoding"`
	Properties      map[string]any `json:"properties"`
}

func (s *MessagesService) Publish(ctx context.Context, vhost, exchange string, publish PublishRequest) error {
	path := fmt.Sprintf("api/exchanges/%s/%s/publish", url.PathEscape(vhost), url.PathEscape(exchange))
	tflog.Debug(ctx, fmt.Sprintf("service=messages method=Publish path=%s", path))
	_, err := s.client.Request(ctx, http.MethodPost, path, publish)
	return err
}
