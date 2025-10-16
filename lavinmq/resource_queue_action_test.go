package lavinmq

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQueueAction_Purge(t *testing.T) {
	t.Parallel()
	queueActionResourceName := "lavinmq_queue_action.test_action"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "lavinmq_vhost" "test" {
						name = "test_queue_action_purge"
					}

					resource "lavinmq_queue" "test_queue" {
						name    = "test_queue_action_purge"
						vhost   = lavinmq_vhost.test.name
						durable = true
						auto_delete = false
					}

					resource "lavinmq_queue_action" "test_action" {
						name   = lavinmq_queue.test_queue.name
						vhost  = lavinmq_vhost.test.name
						action = "purge"
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(queueActionResourceName, "name", "test_queue_action_purge"),
					resource.TestCheckResourceAttr(queueActionResourceName, "vhost", "test_queue_action_purge"),
					resource.TestCheckResourceAttr(queueActionResourceName, "action", "purge"),
				),
			},
		},
	})
}

func TestAccQueueAction_InvalidAction(t *testing.T) {
	t.Parallel()

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "lavinmq_vhost" "test" {
						name = "test_queue_action_invalid"
					}

					resource "lavinmq_queue" "test_queue" {
						name    = "test_queue_action_invalid"
						vhost   = lavinmq_vhost.test.name
						durable = true
            auto_delete = false
          }
						
					resource "lavinmq_queue_action" "test_action" {
						name   = lavinmq_queue.test_queue.name
						vhost  = lavinmq_vhost.test.name
						action = "invalid_action"
					}
					`,
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}
