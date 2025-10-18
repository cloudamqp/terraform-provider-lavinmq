package lavinmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFederationUpstream_Import(t *testing.T) {
	t.Parallel()
	resourceName := "lavinmq_federation_upstream.test"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "lavinmq_federation_upstream" "test" {
  name     = "vcr_test_federation_upstream_import"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstream-host:5672/%2f"
  exchange = "upstream-exchange"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vcr_test_federation_upstream_import"),
					resource.TestCheckResourceAttr(resourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(resourceName, "uri", "amqp://guest:guest@upstream-host:5672/%2f"),
					resource.TestCheckResourceAttr(resourceName, "exchange", "upstream-exchange"),
					resource.TestCheckResourceAttr(resourceName, "prefetch_count", "1000"),
					resource.TestCheckResourceAttr(resourceName, "ack_mode", "on-confirm"),
				),
			},
			{
				ResourceName:                         resourceName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateId:                        "/@vcr_test_federation_upstream_import",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

func TestAccFederationUpstream_Update(t *testing.T) {
	t.Parallel()
	resourceName := "lavinmq_federation_upstream.test"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "lavinmq_federation_upstream" "test" {
  name           = "vcr_test_federation_upstream_update"
  vhost          = "/"
  uri            = "amqp://guest:guest@upstream-host:5672/%2f"
  exchange       = "upstream-exchange"
  prefetch_count = 500
  ack_mode       = "on-publish"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vcr_test_federation_upstream_update"),
					resource.TestCheckResourceAttr(resourceName, "prefetch_count", "500"),
					resource.TestCheckResourceAttr(resourceName, "ack_mode", "on-publish"),
				),
			},
			{
				Config: `
resource "lavinmq_federation_upstream" "test" {
  name            = "vcr_test_federation_upstream_update"
  vhost           = "/"
  uri             = "amqp://guest:guest@upstream-host:5672/%2f"
  exchange        = "upstream-exchange"
  prefetch_count  = 2000
  ack_mode        = "no-ack"
  reconnect_delay = 10
  max_hops        = 3
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vcr_test_federation_upstream_update"),
					resource.TestCheckResourceAttr(resourceName, "prefetch_count", "2000"),
					resource.TestCheckResourceAttr(resourceName, "ack_mode", "no-ack"),
					resource.TestCheckResourceAttr(resourceName, "reconnect_delay", "10"),
					resource.TestCheckResourceAttr(resourceName, "max_hops", "3"),
				),
			},
		},
	})
}

func TestAccFederationUpstream_ExchangeFederation(t *testing.T) {
	t.Parallel()
	resourceName := "lavinmq_federation_upstream.test_exchange"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "lavinmq_federation_upstream" "test_exchange" {
  name     = "vcr_test_exchange_federation"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstream-host:5672/%2f"
  exchange = "federated-exchange"
  max_hops = 2
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vcr_test_exchange_federation"),
					resource.TestCheckResourceAttr(resourceName, "exchange", "federated-exchange"),
					resource.TestCheckResourceAttr(resourceName, "max_hops", "2"),
					resource.TestCheckNoResourceAttr(resourceName, "queue"),
				),
			},
		},
	})
}

func TestAccFederationUpstream_QueueFederation(t *testing.T) {
	t.Parallel()
	resourceName := "lavinmq_federation_upstream.test_queue"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "lavinmq_federation_upstream" "test_queue" {
  name  = "vcr_test_queue_federation"
  vhost = "/"
  uri   = "amqp://guest:guest@upstream-host:5672/%2f"
  queue = "federated-queue"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vcr_test_queue_federation"),
					resource.TestCheckResourceAttr(resourceName, "queue", "federated-queue"),
					resource.TestCheckNoResourceAttr(resourceName, "exchange"),
				),
			},
		},
	})
}

func TestAccFederationUpstream_WithTTL(t *testing.T) {
	t.Parallel()
	resourceName := "lavinmq_federation_upstream.test_ttl"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "lavinmq_federation_upstream" "test_ttl" {
  name        = "vcr_test_federation_ttl"
  vhost       = "/"
  uri         = "amqp://guest:guest@upstream-host:5672/%2f"
  exchange    = "ttl-exchange"
  expires     = 3600000
  message_ttl = 60000
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vcr_test_federation_ttl"),
					resource.TestCheckResourceAttr(resourceName, "expires", "3600000"),
					resource.TestCheckResourceAttr(resourceName, "message_ttl", "60000"),
				),
			},
		},
	})
}
