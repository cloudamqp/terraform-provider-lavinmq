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
					resource "lavinmq_vhost" "test" {
						name = "test-publish-message-vhost"
					}

					resource "lavinmq_exchange" "topic_exchange" {
						name        = "topic-exchange"
						vhost       = lavinmq_vhost.test.name
						type        = "topic"
						auto_delete = false
						durable     = true
					}
					
					resource "lavinmq_queue" "publish_queue" {
						name        = "publish-message-queue"
						vhost       = lavinmq_vhost.test.name
						durable     = true
						auto_delete = false
					}

					resource "lavinmq_binding" "publish_binding" {
						vhost            = lavinmq_vhost.test.name
						source           = lavinmq_exchange.topic_exchange.name
						destination      = lavinmq_queue.publish_queue.name
						destination_type = "queue"
						routing_key      = "publish.*"
					}

					resource "lavinmq_publish_message" "example_message" {
						vhost       = lavinmq_vhost.test.name
						exchange    = lavinmq_exchange.topic_exchange.name
						routing_key = "publish.test"
						payload     = "1"

						depends_on = [
							lavinmq_binding.publish_binding
						]
					}

					data "lavinmq_queues" "all_queues" {
						vhost = lavinmq_vhost.test.name

						depends_on = [
							lavinmq_publish_message.example_message
						]
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(queueDataSourceName, "vhost", "test-publish-message-vhost"),
					resource.TestCheckResourceAttrSet(queueDataSourceName, "queues.#"),
					resource.TestCheckResourceAttr(queueDataSourceName, "queues.0.ready", "1"),
				),
			},
			{
				Config: `
					resource "lavinmq_vhost" "test" {
						name = "test-publish-message-vhost"
					}

					resource "lavinmq_exchange" "topic_exchange" {
						name        = "topic-exchange"
						vhost       = lavinmq_vhost.test.name
						type        = "topic"
						auto_delete = false
						durable     = true
					}
					
					resource "lavinmq_queue" "publish_queue" {
						name        = "publish-message-queue"
						vhost       = lavinmq_vhost.test.name
						durable     = true
						auto_delete = false
					}

					resource "lavinmq_binding" "publish_binding" {
						vhost            = lavinmq_vhost.test.name
						source           = lavinmq_exchange.topic_exchange.name
						destination      = lavinmq_queue.publish_queue.name
						destination_type = "queue"
						routing_key      = "publish.*"
					}

					resource "lavinmq_publish_message" "example_message" {
						vhost       = lavinmq_vhost.test.name
						exchange    = lavinmq_exchange.topic_exchange.name
						routing_key = "publish.test"
						payload     = "2"

						depends_on = [
							lavinmq_binding.publish_binding
						]
					}

					data "lavinmq_queues" "all_queues" {
						vhost = lavinmq_vhost.test.name

						depends_on = [
							lavinmq_publish_message.example_message
						]
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(queueDataSourceName, "vhost", "test-publish-message-vhost"),
					resource.TestCheckResourceAttrSet(queueDataSourceName, "queues.#"),
					resource.TestCheckResourceAttr(queueDataSourceName, "queues.0.ready", "2"),
				),
			},
		},
	})
}
