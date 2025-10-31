package lavinmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPublishMessage_Publish(t *testing.T) {
	t.Parallel()
	queueDataSourceName := "data.lavinmq_queues.all_queues"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "lavinmq_queue" "publish_queue" {
						name        = "test-publish-message-queue"
						vhost       = "/"
						durable     = true
						auto_delete = false
					}

					data "lavinmq_queues" "all_queues" {
						vhost = "/"

						depends_on = [
							lavinmq_queue.publish_queue
						]
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(queueDataSourceName, "vhost", "/"),
					resource.TestCheckResourceAttrSet(queueDataSourceName, "queues.#"),
					resource.TestCheckTypeSetElemNestedAttrs(queueDataSourceName, "queues.*", map[string]string{
						"name":  "test-publish-message-queue",
						"ready": "0",
					}),
				),
			},
			{
				Config: `
					resource "lavinmq_queue" "publish_queue" {
						name        = "test-publish-message-queue"
						vhost       = "/"
						durable     = true
						auto_delete = false
					}

					resource "lavinmq_publish_message" "example_message" {
						vhost       = "/"
						exchange    = "amq.default"
						routing_key = lavinmq_queue.publish_queue.name
						payload     = "{\"message\": \"VCR test publish\"}"
						properties = {
								content_type = "application/json"
						}
						publish_message_counter = 1
					}

					data "lavinmq_queues" "all_queues" {
						vhost = "/"

						depends_on = [
							lavinmq_publish_message.example_message
						]
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(queueDataSourceName, "vhost", "/"),
					resource.TestCheckResourceAttrSet(queueDataSourceName, "queues.#"),
					resource.TestCheckTypeSetElemNestedAttrs(queueDataSourceName, "queues.*", map[string]string{
						"name":  "test-publish-message-queue",
						"ready": "1",
					}),
				),
			},
			{
				Config: `
					resource "lavinmq_queue" "publish_queue" {
						name        = "test-publish-message-queue"
						vhost       = "/"
						durable     = true
						auto_delete = false
					}

					resource "lavinmq_publish_message" "example_message" {
						vhost       = "/"
						exchange    = "amq.default"
						routing_key = lavinmq_queue.publish_queue.name
						payload     = "{\"message\": \"VCR test publish\"}"
						properties = {
								content_type = "application/json"
						}
						publish_message_counter = 2
					}

					data "lavinmq_queues" "all_queues" {
						vhost = "/"

						depends_on = [
							lavinmq_publish_message.example_message
						]
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(queueDataSourceName, "vhost", "/"),
					resource.TestCheckResourceAttrSet(queueDataSourceName, "queues.#"),
					resource.TestCheckTypeSetElemNestedAttrs(queueDataSourceName, "queues.*", map[string]string{
						"name":  "test-publish-message-queue",
						"ready": "2",
					}),
				),
			},
		},
	})
}
