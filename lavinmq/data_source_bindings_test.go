package lavinmq

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceBindings_Basic(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.lavinmq_bindings.test"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_exchange" "test_exchange" {
            name        = "vcr_test_exchange_for_data_source"
            vhost       = "/"
            type        = "direct"
            durable     = true
            auto_delete = false
          }

          resource "lavinmq_queue" "test_queue" {
            name        = "vcr_test_queue_for_data_source"
            vhost       = "/"
            durable     = true
            auto_delete = false
          }

          resource "lavinmq_binding" "test_binding" {
            vhost            = "/"
            source           = lavinmq_exchange.test_exchange.name
            destination      = lavinmq_queue.test_queue.name
            destination_type = "queue"
            routing_key      = "test.datasource"
          }

          data "lavinmq_bindings" "test" {
            vhost = "/"
            depends_on = [lavinmq_binding.test_binding]
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "vhost", "/"),
					resource.TestCheckResourceAttrSet(dataSourceName, "bindings.#"),
				),
			},
		},
	})
}

func TestAccDataSourceBindings_NonExistingVhost(t *testing.T) {
	t.Parallel()
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          data "lavinmq_bindings" "empty" {
            vhost = "terraform-lavinmq-non-existing-bindings-test"
          }`,
				ExpectError: regexp.MustCompile(`status code: 403|Access Refused`),
			},
		},
	})
}
