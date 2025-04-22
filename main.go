package main

import (
	"context"

	"github.com/cloudamqp/terraform-provider-lavinmq/lavinmq"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), lavinmq.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/cloudamqp/lavinmq",
	})
}
