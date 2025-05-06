package main

import (
	"context"
	"net/http"

	"github.com/cloudamqp/terraform-provider-lavinmq/lavinmq"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	err := providerserver.Serve(context.Background(), func() provider.Provider {
		return lavinmq.New("1.0", http.DefaultClient)
	}, providerserver.ServeOpts{
		Address: "registry.terraform.io/cloudamqp/lavinmq",
	})
	if err != nil {
		panic(err)
	}
}
