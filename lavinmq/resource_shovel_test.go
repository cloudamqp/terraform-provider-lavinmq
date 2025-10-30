package lavinmq

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccShovel_Import(t *testing.T) {
	t.Parallel()
	shovelResourceName := "lavinmq_shovel.test_shovel"

	// Set sanitized value for playback and use real value for recording
	testSrcDestURI := "SHOVEL_SRC_DEST_URI"
	if os.Getenv("LAVINMQ_RECORD") != "" {
		testSrcDestURI = os.Getenv("SHOVEL_SRC_DEST_URI")
	}

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_shovel" "test_shovel" {
						name        = "vcr_test_shovel_import"
						vhost       = "/"
						src_uri     = "%[1]s"
						dest_uri    = "%[1]s"
						src_queue   = "source_queue"
						dest_queue  = "dest_queue"
					}`, testSrcDestURI),
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

	// Set sanitized value for playback and use real value for recording
	testSrcDestURI := "SHOVEL_SRC_DEST_URI"
	if os.Getenv("LAVINMQ_RECORD") != "" {
		testSrcDestURI = os.Getenv("SHOVEL_SRC_DEST_URI")
	}

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_shovel" "test_shovel" {
						name               = "vcr_test_shovel_update"
						vhost              = "/"
						src_uri            = "%[1]s"
						dest_uri           = "%[1]s"
						src_queue          = "source_queue"
						dest_queue         = "dest_queue"
						src_prefetch_count = 500
						ack_mode           = "on-publish"
					}`, testSrcDestURI),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(shovelResourceName, "name", "vcr_test_shovel_update"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_prefetch_count", "500"),
					resource.TestCheckResourceAttr(shovelResourceName, "ack_mode", "on-publish"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_shovel" "test_shovel" {
						name               = "vcr_test_shovel_update"
						vhost              = "/"
						src_uri            = "%[1]s"
						dest_uri           = "%[1]s"
						src_queue          = "source_queue"
						dest_queue         = "dest_queue"
						src_prefetch_count = 2000
						ack_mode           = "no-ack"
						reconnect_delay    = 10
					}`, testSrcDestURI),
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

	// Set sanitized value for playback and use real value for recording
	testSrcDestURI := "SHOVEL_SRC_DEST_URI"
	if os.Getenv("LAVINMQ_RECORD") != "" {
		testSrcDestURI = os.Getenv("SHOVEL_SRC_DEST_URI")
	}

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_shovel" "test_q2q" {
						name       = "vcr_test_q2q"
						vhost      = "/"
						src_uri    = "%[1]s"
						dest_uri   = "%[1]s"
						src_queue  = "source_queue"
						dest_queue = "dest_queue"
					}`, testSrcDestURI),
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

	// Set sanitized value for playback and use real value for recording
	testSrcDestURI := "SHOVEL_SRC_DEST_URI"
	if os.Getenv("LAVINMQ_RECORD") != "" {
		testSrcDestURI = os.Getenv("SHOVEL_SRC_DEST_URI")
	}

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_shovel" "test_q2e" {
						name             = "vcr_test_q2e"
						vhost            = "/"
						src_uri          = "%[1]s"
						dest_uri         = "%[1]s"
						src_queue        = "source_queue"
						dest_exchange    = "dest_exchange"
						dest_exchange_key = "routing.key"
					}`, testSrcDestURI),
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

	// Set sanitized value for playback and use real value for recording
	testSrcDestURI := "SHOVEL_SRC_DEST_URI"
	if os.Getenv("LAVINMQ_RECORD") != "" {
		testSrcDestURI = os.Getenv("SHOVEL_SRC_DEST_URI")
	}

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_shovel" "test_e2q" {
						name              = "vcr_test_e2q"
						vhost             = "/"
						src_uri           = "%[1]s"
						dest_uri          = "%[1]s"
						src_exchange      = "source_exchange"
						src_exchange_key  = "source.key"
						dest_queue        = "dest_queue"
					}`, testSrcDestURI),
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

	// Set sanitized value for playback and use real value for recording
	testSrcDestURI := "SHOVEL_SRC_DEST_URI"
	if os.Getenv("LAVINMQ_RECORD") != "" {
		testSrcDestURI = os.Getenv("SHOVEL_SRC_DEST_URI")
	}

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_shovel" "test_e2e" {
						name              = "vcr_test_e2e"
						vhost             = "/"
						src_uri           = "%[1]s"
						dest_uri          = "%[1]s"
						src_exchange      = "source_exchange"
						src_exchange_key  = "source.key"
						dest_exchange     = "dest_exchange"
						dest_exchange_key = "dest.key"
					}`, testSrcDestURI),
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

	// Set sanitized value for playback and use real value for recording
	testSrcDestURI := "SHOVEL_SRC_DEST_URI"
	if os.Getenv("LAVINMQ_RECORD") != "" {
		testSrcDestURI = os.Getenv("SHOVEL_SRC_DEST_URI")
	}

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_shovel" "test_invalid" {
						name         = "vcr_test_invalid"
						vhost        = "/"
						src_uri      = "%[1]s"
						dest_uri     = "%[1]s"
						src_queue    = "source_queue"
						src_exchange = "source_exchange"
						dest_queue   = "dest_queue"
					}`, testSrcDestURI),
				ExpectError: regexp.MustCompile(`Cannot specify both src_queue and src_exchange`),
			},
		},
	})
}

func TestAccShovel_InvalidBothDestinations(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testSrcDestURI := "SHOVEL_SRC_DEST_URI"
	if os.Getenv("LAVINMQ_RECORD") != "" {
		testSrcDestURI = os.Getenv("SHOVEL_SRC_DEST_URI")
	}

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_shovel" "test_invalid" {
						name          = "vcr_test_invalid"
						vhost         = "/"
						src_uri       = "%[1]s"
						dest_uri      = "%[1]s"
						src_queue     = "source_queue"
						dest_queue    = "dest_queue"
						dest_exchange = "dest_exchange"
					}`, testSrcDestURI),
				ExpectError: regexp.MustCompile(`Cannot specify both dest_queue and dest_exchange`),
			},
		},
	})
}
