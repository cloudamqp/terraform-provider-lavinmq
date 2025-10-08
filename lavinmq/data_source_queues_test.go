package lavinmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceQueues_Basic(t *testing.T) {
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          data "lavinmq_queues" "all" {
            vhost = "/"
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lavinmq_queues.all", "vhost", "/"),
					resource.TestCheckResourceAttrSet("data.lavinmq_queues.all", "queues.#"),
				),
			},
		},
	})
}

func TestAccDataSourceQueues_Empty(t *testing.T) {
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_vhost" "test" {
            name = "terraform-lavinmq-empty-test"
          }

          data "lavinmq_queues" "empty" {
            vhost = lavinmq_vhost.test.name
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lavinmq_queues.empty", "vhost", "terraform-lavinmq-empty-test"),
					resource.TestCheckResourceAttr("data.lavinmq_queues.empty", "queues.#", "0"),
				),
			},
		},
	})
}
