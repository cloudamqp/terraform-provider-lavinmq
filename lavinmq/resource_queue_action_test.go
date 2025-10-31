package lavinmq

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQueueAction_Purge(t *testing.T) {
	t.Parallel()
	queueActionResourceName := "lavinmq_queue_action.test_action"
	queueDataSourceName := "data.lavinmq_queues.all_queues"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "lavinmq_queue" "purge_queue" {
						name    = "test-purge-queue"
						vhost   = "/"
						durable = true
						auto_delete = false
					}

					data "lavinmq_queues" "all_queues" {
						vhost = "/"

						depends_on = [
							lavinmq_queue.purge_queue
						]
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(queueDataSourceName, "vhost", "/"),
					resource.TestCheckResourceAttrSet(queueDataSourceName, "queues.#"),
					resource.TestCheckTypeSetElemNestedAttrs(queueDataSourceName, "queues.*", map[string]string{
						"name":  "test-purge-queue",
						"ready": "0",
					}),
				),
			},
			{
				Config: `
					resource "lavinmq_queue" "purge_queue" {
						name    = "test-purge-queue"
						vhost   = "/"
						durable = true
						auto_delete = false
					}

					resource "lavinmq_publish_message" "test_message" {
						vhost       = "/"
						exchange    = "amq.default"
						routing_key = lavinmq_queue.purge_queue.name
						payload     = "{\"message\": \"VCR test publish\"}"
						properties = {
								content_type = "application/json"
						}
						publish_message_counter = 1

						depends_on = [
							lavinmq_queue.purge_queue
						]
					}

					data "lavinmq_queues" "all_queues" {
						vhost = "/"

						depends_on = [
							lavinmq_publish_message.test_message
						]
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(queueDataSourceName, "vhost", "/"),
					resource.TestCheckResourceAttrSet(queueDataSourceName, "queues.#"),
					resource.TestCheckTypeSetElemNestedAttrs(queueDataSourceName, "queues.*", map[string]string{
						"name":  "test-purge-queue",
						"ready": "1",
					}),
				),
			},
			{
				Config: `
					resource "lavinmq_queue" "purge_queue" {
						name    = "test-purge-queue"
						vhost   = "/"
						durable = true
						auto_delete = false
					}

					resource "lavinmq_publish_message" "test_message" {
						vhost       = "/"
						exchange    = "amq.default"
						routing_key = lavinmq_queue.purge_queue.name
						payload     = "{\"message\": \"VCR test publish\"}"
						properties = {
								content_type = "application/json"
						}
						publish_message_counter = 2

						depends_on = [
							lavinmq_queue.purge_queue
						]
					}

					resource "lavinmq_queue_action" "test_action" {
						name   = lavinmq_queue.purge_queue.name
						vhost  = "/"
						action = "purge"

						depends_on = [
							lavinmq_publish_message.test_message
						]
					}

					data "lavinmq_queues" "all_queues" {
						vhost = "/"

						depends_on = [
							lavinmq_queue_action.test_action
						]
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(queueActionResourceName, "name", "test-purge-queue"),
					resource.TestCheckResourceAttr(queueActionResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(queueActionResourceName, "action", "purge"),
					resource.TestCheckResourceAttr(queueDataSourceName, "vhost", "/"),
					resource.TestCheckResourceAttrSet(queueDataSourceName, "queues.#"),
					resource.TestCheckTypeSetElemNestedAttrs(queueDataSourceName, "queues.*", map[string]string{
						"name":  "test-purge-queue",
						"ready": "0",
					}),
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
					resource "lavinmq_queue" "test_queue" {
						name    = "test_queue_action_invalid"
						vhost   = "/"
						durable = true
            auto_delete = false
          }

					resource "lavinmq_queue_action" "test_action" {
						name   = lavinmq_queue.test_queue.name
						vhost  = "/"
						action = "invalid_action"
					}`,
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}
