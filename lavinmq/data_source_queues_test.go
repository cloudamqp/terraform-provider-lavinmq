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
          resource "lavinmq_vhost" "test" {
            name = "terraform-lavinmq-test"
          }

          resource "lavinmq_queue" "test1" {
            name       = "terraform-queue-test-1"
            vhost      = lavinmq_vhost.test.name
            durable    = true
            auto_delete = false
          }

          resource "lavinmq_queue" "test2" {
            name       = "terraform-queue-test-2"
            vhost      = lavinmq_vhost.test.name
            durable    = false
            auto_delete = true
          }

          data "lavinmq_queues" "all" {
            vhost = lavinmq_vhost.test.name
            depends_on = [lavinmq_queue.test1, lavinmq_queue.test2]
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lavinmq_queues.all", "vhost", "terraform-lavinmq-test"),
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_queues.all", "queues.*", map[string]string{
						"name":        "terraform-queue-test-1",
						"vhost":       "terraform-lavinmq-test",
						"durable":     "true",
						"auto_delete": "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_queues.all", "queues.*", map[string]string{
						"name":        "terraform-queue-test-2",
						"vhost":       "terraform-lavinmq-test",
						"durable":     "false",
						"auto_delete": "true",
					}),
				),
			},
		},
	})
}
