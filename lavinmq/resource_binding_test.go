package lavinmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBinding_Basic(t *testing.T) {
	t.Parallel()
	bindingResourceName := "lavinmq_binding.test_binding"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_exchange" "test_exchange" {
            name        = "vcr_test_exchange_for_binding"
            vhost       = "/"
            type        = "direct"
            durable     = true
            auto_delete = false
          }

          resource "lavinmq_queue" "test_queue" {
            name        = "vcr_test_queue_for_binding"
            vhost       = "/"
            durable     = true
            auto_delete = false
          }

          resource "lavinmq_binding" "test_binding" {
            vhost            = "/"
            source           = lavinmq_exchange.test_exchange.name
            destination      = lavinmq_queue.test_queue.name
            destination_type = "queue"
            routing_key      = "test.key"
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(bindingResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(bindingResourceName, "source", "vcr_test_exchange_for_binding"),
					resource.TestCheckResourceAttr(bindingResourceName, "destination", "vcr_test_queue_for_binding"),
					resource.TestCheckResourceAttr(bindingResourceName, "destination_type", "queue"),
					resource.TestCheckResourceAttr(bindingResourceName, "routing_key", "test.key"),
					resource.TestCheckResourceAttrSet(bindingResourceName, "properties_key"),
				),
			},
		},
	})
}

func TestAccBinding_Import(t *testing.T) {
	t.Parallel()
	bindingResourceName := "lavinmq_binding.test_binding"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_exchange" "test_exchange" {
            name        = "vcr_test_exchange_for_binding_import"
            vhost       = "/"
            type        = "direct"
            durable     = true
            auto_delete = false
          }

          resource "lavinmq_queue" "test_queue" {
            name        = "vcr_test_queue_for_binding_import"
            vhost       = "/"
            durable     = true
            auto_delete = false
          }

          resource "lavinmq_binding" "test_binding" {
            vhost            = "/"
            source           = lavinmq_exchange.test_exchange.name
            destination      = lavinmq_queue.test_queue.name
            destination_type = "queue"
            routing_key      = "import.key"
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(bindingResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(bindingResourceName, "source", "vcr_test_exchange_for_binding_import"),
					resource.TestCheckResourceAttr(bindingResourceName, "destination", "vcr_test_queue_for_binding_import"),
					resource.TestCheckResourceAttr(bindingResourceName, "destination_type", "queue"),
					resource.TestCheckResourceAttr(bindingResourceName, "routing_key", "import.key"),
				),
			},
			{
				ResourceName:                         bindingResourceName,
				ImportStateVerifyIdentifierAttribute: "properties_key",
				ImportStateId:                        "/@vcr_test_exchange_for_binding_import@vcr_test_queue_for_binding_import@queue@import.key",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

func TestAccBinding_HeadersExchange(t *testing.T) {
	t.Parallel()
	bindingResourceName := "lavinmq_binding.test_binding"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_exchange" "test_exchange" {
            name        = "vcr_test_headers_exchange"
            vhost       = "/"
            type        = "headers"
            durable     = true
            auto_delete = false
          }

          resource "lavinmq_queue" "test_queue" {
            name        = "vcr_test_queue_headers"
            vhost       = "/"
            durable     = true
            auto_delete = false
          }

          resource "lavinmq_binding" "test_binding" {
            vhost            = "/"
            source           = lavinmq_exchange.test_exchange.name
            destination      = lavinmq_queue.test_queue.name
            destination_type = "queue"
            routing_key      = ""
            arguments = {
              x-match  = "all"
              priority = "high"
              type     = "alert"
            }
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(bindingResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(bindingResourceName, "source", "vcr_test_headers_exchange"),
					resource.TestCheckResourceAttr(bindingResourceName, "destination", "vcr_test_queue_headers"),
					resource.TestCheckResourceAttr(bindingResourceName, "destination_type", "queue"),
					resource.TestCheckResourceAttr(bindingResourceName, "routing_key", ""),
					resource.TestCheckResourceAttr(bindingResourceName, "arguments.x-match", "all"),
					resource.TestCheckResourceAttr(bindingResourceName, "arguments.priority", "high"),
					resource.TestCheckResourceAttr(bindingResourceName, "arguments.type", "alert"),
					resource.TestCheckResourceAttrSet(bindingResourceName, "properties_key"),
				),
			},
		},
	})
}

func TestAccBinding_TopicExchange(t *testing.T) {
	t.Parallel()
	bindingResourceName := "lavinmq_binding.test_binding"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_exchange" "test_exchange" {
            name        = "vcr_test_topic_exchange"
            vhost       = "/"
            type        = "topic"
            durable     = true
            auto_delete = false
          }

          resource "lavinmq_queue" "test_queue" {
            name        = "vcr_test_queue_topic"
            vhost       = "/"
            durable     = true
            auto_delete = false
          }

          resource "lavinmq_binding" "test_binding" {
            vhost            = "/"
            source           = lavinmq_exchange.test_exchange.name
            destination      = lavinmq_queue.test_queue.name
            destination_type = "queue"
            routing_key      = "events.#"
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(bindingResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(bindingResourceName, "source", "vcr_test_topic_exchange"),
					resource.TestCheckResourceAttr(bindingResourceName, "destination", "vcr_test_queue_topic"),
					resource.TestCheckResourceAttr(bindingResourceName, "destination_type", "queue"),
					resource.TestCheckResourceAttr(bindingResourceName, "routing_key", "events.#"),
					resource.TestCheckResourceAttrSet(bindingResourceName, "properties_key"),
				),
			},
		},
	})
}

func TestAccBinding_FanoutExchange(t *testing.T) {
	t.Parallel()
	bindingResourceName := "lavinmq_binding.test_binding"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_exchange" "test_exchange" {
            name        = "vcr_test_fanout_exchange"
            vhost       = "/"
            type        = "fanout"
            durable     = true
            auto_delete = false
          }

          resource "lavinmq_queue" "test_queue" {
            name        = "vcr_test_queue_fanout"
            vhost       = "/"
            durable     = true
            auto_delete = false
          }

          resource "lavinmq_binding" "test_binding" {
            vhost            = "/"
            source           = lavinmq_exchange.test_exchange.name
            destination      = lavinmq_queue.test_queue.name
            destination_type = "queue"
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(bindingResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(bindingResourceName, "source", "vcr_test_fanout_exchange"),
					resource.TestCheckResourceAttr(bindingResourceName, "destination", "vcr_test_queue_fanout"),
					resource.TestCheckResourceAttr(bindingResourceName, "destination_type", "queue"),
					resource.TestCheckResourceAttr(bindingResourceName, "routing_key", ""),
					resource.TestCheckResourceAttrSet(bindingResourceName, "properties_key"),
				),
			},
		},
	})
}

func TestAccBinding_DirectExchange(t *testing.T) {
	t.Parallel()
	bindingResourceName := "lavinmq_binding.test_binding"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_exchange" "test_exchange" {
            name        = "vcr_test_direct_exchange"
            vhost       = "/"
            type        = "direct"
            durable     = true
            auto_delete = false
          }

          resource "lavinmq_queue" "test_queue" {
            name        = "vcr_test_queue_direct"
            vhost       = "/"
            durable     = true
            auto_delete = false
          }

          resource "lavinmq_binding" "test_binding" {
            vhost            = "/"
            source           = lavinmq_exchange.test_exchange.name
            destination      = lavinmq_queue.test_queue.name
            destination_type = "queue"
            routing_key      = "specific.route"
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(bindingResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(bindingResourceName, "source", "vcr_test_direct_exchange"),
					resource.TestCheckResourceAttr(bindingResourceName, "destination", "vcr_test_queue_direct"),
					resource.TestCheckResourceAttr(bindingResourceName, "destination_type", "queue"),
					resource.TestCheckResourceAttr(bindingResourceName, "routing_key", "specific.route"),
					resource.TestCheckResourceAttrSet(bindingResourceName, "properties_key"),
				),
			},
		},
	})
}

func TestAccBinding_ExchangeToExchange(t *testing.T) {
	t.Parallel()
	bindingResourceName := "lavinmq_binding.test_binding"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_exchange" "source_exchange" {
            name        = "vcr_test_source_exchange"
            vhost       = "/"
            type        = "topic"
            durable     = true
            auto_delete = false
          }

          resource "lavinmq_exchange" "dest_exchange" {
            name        = "vcr_test_dest_exchange"
            vhost       = "/"
            type        = "fanout"
            durable     = true
            auto_delete = false
          }

          resource "lavinmq_binding" "test_binding" {
            vhost            = "/"
            source           = lavinmq_exchange.source_exchange.name
            destination      = lavinmq_exchange.dest_exchange.name
            destination_type = "exchange"
            routing_key      = "backup.#"
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(bindingResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(bindingResourceName, "source", "vcr_test_source_exchange"),
					resource.TestCheckResourceAttr(bindingResourceName, "destination", "vcr_test_dest_exchange"),
					resource.TestCheckResourceAttr(bindingResourceName, "destination_type", "exchange"),
					resource.TestCheckResourceAttr(bindingResourceName, "routing_key", "backup.#"),
					resource.TestCheckResourceAttrSet(bindingResourceName, "properties_key"),
				),
			},
		},
	})
}
