package lavinmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceExchanges_Basic(t *testing.T) {
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_exchange" "test1" {
            name    = "terraform-lavinmq-test-exchange-1"
            vhost   = "/"
            type    = "direct"
            durable = true
          }

          resource "lavinmq_exchange" "test2" {
            name        = "terraform-lavinmq-test-exchange-2"
            vhost       = "/"
            type        = "fanout"
            durable     = false
            auto_delete = true
          }

          data "lavinmq_exchanges" "all" {
            depends_on = [lavinmq_exchange.test1, lavinmq_exchange.test2]
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_exchanges.all", "exchanges.*", map[string]string{
						"name":        "terraform-lavinmq-test-exchange-1",
						"vhost":       "/",
						"type":        "direct",
						"durable":     "true",
						"auto_delete": "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_exchanges.all", "exchanges.*", map[string]string{
						"name":        "terraform-lavinmq-test-exchange-2",
						"vhost":       "/",
						"type":        "fanout",
						"durable":     "false",
						"auto_delete": "true",
					}),
				),
			},
		},
	})
}
