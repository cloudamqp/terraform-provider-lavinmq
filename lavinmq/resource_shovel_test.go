package lavinmq

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccShovel_Import(t *testing.T) {
	t.Parallel()
	shovelResourceName := "lavinmq_shovel.test_shovel"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "lavinmq_shovel" "test_shovel" {
  name        = "vcr_test_shovel_import"
  vhost       = "/"
  src_uri     = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri    = "amqp://guest:guest@localhost:5672/%2f"
  src_queue   = "source_queue"
  dest_queue  = "dest_queue"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(shovelResourceName, "name", "vcr_test_shovel_import"),
					resource.TestCheckResourceAttr(shovelResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_queue", "source_queue"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_queue", "dest_queue"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_prefetch_count", "1000"),
					resource.TestCheckResourceAttr(shovelResourceName, "ack_mode", "on-confirm"),
				),
			},
			{
				ResourceName:                         shovelResourceName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateId:                        "/@vcr_test_shovel_import",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

func TestAccShovel_Update(t *testing.T) {
	t.Parallel()
	shovelResourceName := "lavinmq_shovel.test_shovel"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "lavinmq_shovel" "test_shovel" {
  name               = "vcr_test_shovel_update"
  vhost              = "/"
  src_uri            = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri           = "amqp://guest:guest@localhost:5672/%2f"
  src_queue          = "source_queue"
  dest_queue         = "dest_queue"
  src_prefetch_count = 500
  ack_mode           = "on-publish"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(shovelResourceName, "name", "vcr_test_shovel_update"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_prefetch_count", "500"),
					resource.TestCheckResourceAttr(shovelResourceName, "ack_mode", "on-publish"),
				),
			},
			{
				Config: `
resource "lavinmq_shovel" "test_shovel" {
  name               = "vcr_test_shovel_update"
  vhost              = "/"
  src_uri            = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri           = "amqp://guest:guest@localhost:5672/%2f"
  src_queue          = "source_queue"
  dest_queue         = "dest_queue"
  src_prefetch_count = 2000
  ack_mode           = "no-ack"
  reconnect_delay    = 10
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(shovelResourceName, "name", "vcr_test_shovel_update"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_prefetch_count", "2000"),
					resource.TestCheckResourceAttr(shovelResourceName, "ack_mode", "no-ack"),
					resource.TestCheckResourceAttr(shovelResourceName, "reconnect_delay", "10"),
				),
			},
		},
	})
}

func TestAccShovel_QueueToQueue(t *testing.T) {
	t.Parallel()
	shovelResourceName := "lavinmq_shovel.test_q2q"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "lavinmq_shovel" "test_q2q" {
  name       = "vcr_test_q2q"
  vhost      = "/"
  src_uri    = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri   = "amqp://guest:guest@localhost:5672/%2f"
  src_queue  = "source_queue"
  dest_queue = "dest_queue"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(shovelResourceName, "name", "vcr_test_q2q"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_queue", "source_queue"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_queue", "dest_queue"),
					resource.TestCheckNoResourceAttr(shovelResourceName, "src_exchange"),
					resource.TestCheckNoResourceAttr(shovelResourceName, "dest_exchange"),
				),
			},
		},
	})
}

func TestAccShovel_QueueToExchange(t *testing.T) {
	t.Parallel()
	shovelResourceName := "lavinmq_shovel.test_q2e"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "lavinmq_shovel" "test_q2e" {
  name             = "vcr_test_q2e"
  vhost            = "/"
  src_uri          = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri         = "amqp://guest:guest@localhost:5672/%2f"
  src_queue        = "source_queue"
  dest_exchange    = "dest_exchange"
  dest_exchange_key = "routing.key"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(shovelResourceName, "name", "vcr_test_q2e"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_queue", "source_queue"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_exchange", "dest_exchange"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_exchange_key", "routing.key"),
				),
			},
		},
	})
}

func TestAccShovel_ExchangeToQueue(t *testing.T) {
	t.Parallel()
	shovelResourceName := "lavinmq_shovel.test_e2q"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "lavinmq_shovel" "test_e2q" {
  name              = "vcr_test_e2q"
  vhost             = "/"
  src_uri           = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri          = "amqp://guest:guest@localhost:5672/%2f"
  src_exchange      = "source_exchange"
  src_exchange_key  = "source.key"
  dest_queue        = "dest_queue"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(shovelResourceName, "name", "vcr_test_e2q"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_exchange", "source_exchange"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_exchange_key", "source.key"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_queue", "dest_queue"),
				),
			},
		},
	})
}

func TestAccShovel_ExchangeToExchange(t *testing.T) {
	t.Parallel()
	shovelResourceName := "lavinmq_shovel.test_e2e"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "lavinmq_shovel" "test_e2e" {
  name              = "vcr_test_e2e"
  vhost             = "/"
  src_uri           = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri          = "amqp://guest:guest@localhost:5672/%2f"
  src_exchange      = "source_exchange"
  src_exchange_key  = "source.key"
  dest_exchange     = "dest_exchange"
  dest_exchange_key = "dest.key"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(shovelResourceName, "name", "vcr_test_e2e"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_exchange", "source_exchange"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_exchange_key", "source.key"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_exchange", "dest_exchange"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_exchange_key", "dest.key"),
				),
			},
		},
	})
}

func TestAccShovel_InvalidBothSources(t *testing.T) {
	t.Parallel()
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "lavinmq_shovel" "test_invalid" {
  name         = "vcr_test_invalid"
  vhost        = "/"
  src_uri      = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri     = "amqp://guest:guest@localhost:5672/%2f"
  src_queue    = "source_queue"
  src_exchange = "source_exchange"
  dest_queue   = "dest_queue"
}`,
				ExpectError: regexp.MustCompile(`Cannot specify both src_queue and src_exchange`),
			},
		},
	})
}

func TestAccShovel_InvalidBothDestinations(t *testing.T) {
	t.Parallel()
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "lavinmq_shovel" "test_invalid" {
  name          = "vcr_test_invalid"
  vhost         = "/"
  src_uri       = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri      = "amqp://guest:guest@localhost:5672/%2f"
  src_queue     = "source_queue"
  dest_queue    = "dest_queue"
  dest_exchange = "dest_exchange"
}`,
				ExpectError: regexp.MustCompile(`Cannot specify both dest_queue and dest_exchange`),
			},
		},
	})
}
