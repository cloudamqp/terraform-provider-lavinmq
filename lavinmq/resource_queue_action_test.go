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

					resource "lavinmq_exchange" "test_exchange" {
						name        = "test_exchange"
						vhost       = lavinmq_vhost.test.name
						type        = "topic"
						auto_delete = false
						durable     = true
					}

					resource "lavinmq_queue" "test_queue" {
						name    = "test_queue_action_purge"
						vhost   = lavinmq_vhost.test.name
						durable = true
						auto_delete = false
					}

					resource "lavinmq_binding" "test_binding" {
						vhost            = lavinmq_vhost.test.name
						source           = lavinmq_exchange.test_exchange.name
						destination      = lavinmq_queue.test_queue.name
						destination_type = "queue"
						routing_key      = "test.routing.key"
					}

					resource "lavinmq_publish_message" "test_message" {
						vhost       = lavinmq_vhost.test.name
						exchange    = lavinmq_exchange.test_exchange.name
						routing_key = "test.routing.key"
						payload     = "Test message for purge action"

						depends_on = [
							lavinmq_binding.test_binding
						]
					}

					data "lavinmq_queues" "all_queues" {
						vhost = lavinmq_vhost.test.name

						depends_on = [
							lavinmq_publish_message.test_message
						]
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lavinmq_queues.all_queues", "vhost", "test_queue_action_purge"),
					resource.TestCheckResourceAttrSet("data.lavinmq_queues.all_queues", "queues.#"),
					resource.TestCheckResourceAttr("data.lavinmq_queues.all_queues", "queues.0.ready", "1"),
				),
			},
			{
				Config: `
					resource "lavinmq_vhost" "test" {
						name = "test_queue_action_purge"
					}

					resource "lavinmq_exchange" "test_exchange" {
						name        = "test_exchange"
						vhost       = lavinmq_vhost.test.name
						type        = "topic"
						auto_delete = false
						durable     = true
					}

					resource "lavinmq_queue" "test_queue" {
						name    = "test_queue_action_purge"
						vhost   = lavinmq_vhost.test.name
						durable = true
						auto_delete = false
					}

					resource "lavinmq_binding" "test_binding" {
						vhost            = lavinmq_vhost.test.name
						source           = lavinmq_exchange.test_exchange.name
						destination      = lavinmq_queue.test_queue.name
						destination_type = "queue"
						routing_key      = "test.routing.key"
					}

					resource "lavinmq_queue_action" "test_action" {
						name   = lavinmq_queue.test_queue.name
						vhost  = lavinmq_vhost.test.name
						action = "purge"

						depends_on = [
							lavinmq_binding.test_binding
						]
					}

					data "lavinmq_queues" "all_queues" {
						vhost = lavinmq_vhost.test.name

						depends_on = [
							lavinmq_queue_action.test_action
						]
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lavinmq_queues.all_queues", "vhost", "test_queue_action_purge"),
					resource.TestCheckResourceAttrSet("data.lavinmq_queues.all_queues", "queues.#"),
					resource.TestCheckResourceAttr(queueActionResourceName, "name", "test_queue_action_purge"),
					resource.TestCheckResourceAttr(queueActionResourceName, "vhost", "test_queue_action_purge"),
					resource.TestCheckResourceAttr("data.lavinmq_queues.all_queues", "queues.0.ready", "0"),
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
