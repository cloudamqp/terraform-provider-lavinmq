package lavinmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceExchanges_Basic(t *testing.T) {
	t.Parallel()
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

func TestAccDataSourceExchanges_DefaultExchanges(t *testing.T) {
	t.Parallel()
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_vhost" "test" {
            name = "terraform-lavinmq-default-exchanges-test"
          }

          data "lavinmq_exchanges" "default" {
            vhost = lavinmq_vhost.test.name
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lavinmq_exchanges.default", "vhost", "terraform-lavinmq-default-exchanges-test"),
					resource.TestCheckResourceAttrSet("data.lavinmq_exchanges.default", "exchanges.#"),
				),
			},
		},
	})
}

func TestAccDataSourceExchanges_NonExistingVhost(t *testing.T) {
	t.Parallel()
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          data "lavinmq_exchanges" "empty" {
            vhost = "terraform-lavinmq-non-existing-test"
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lavinmq_exchanges.empty", "vhost", "terraform-lavinmq-non-existing-test"),
					resource.TestCheckResourceAttr("data.lavinmq_exchanges.empty", "exchanges.#", "0"),
				),
			},
		},
	})
}
