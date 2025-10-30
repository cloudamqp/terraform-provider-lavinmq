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
	queueDataSourceName := "data.lavinmq_queues.all_queues"

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
					resource "lavinmq_queue" "source_queue" {
						name        = "test_q2q_source_queue"
						vhost       = "/"
						durable     = true
						auto_delete = false
					}

					resource "lavinmq_queue" "dest_queue" {
						name        = "test_q2q_dest_queue"
						vhost       = "/"
						durable     = true
						auto_delete = false
					}

					resource "lavinmq_shovel" "test_q2q" {
						name       = "vcr_test_q2q"
						vhost      = "/"
						src_uri    = "%[1]s"
						dest_uri   = "%[1]s"
						src_queue  = lavinmq_queue.source_queue.name
						dest_queue = lavinmq_queue.dest_queue.name
					}

					data "lavinmq_queues" "all_queues" {
						vhost = "/"

						depends_on = [
							lavinmq_shovel.test_q2q
						]
					}`, testSrcDestURI),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(shovelResourceName, "name", "vcr_test_q2q"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_queue", "test_q2q_source_queue"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_queue", "test_q2q_dest_queue"),
					resource.TestCheckNoResourceAttr(shovelResourceName, "src_exchange"),
					resource.TestCheckNoResourceAttr(shovelResourceName, "dest_exchange"),
					resource.TestCheckResourceAttrSet(queueDataSourceName, "queues.#"),
					resource.TestCheckTypeSetElemNestedAttrs(queueDataSourceName, "queues.*", map[string]string{
						"name":  "test_q2q_dest_queue",
						"ready": "0",
					}),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_queue" "source_queue" {
						name        = "test_q2q_source_queue"
						vhost       = "/"
						durable     = true
						auto_delete = false
					}

					resource "lavinmq_queue" "dest_queue" {
						name        = "test_q2q_dest_queue"
						vhost       = "/"
						durable     = true
						auto_delete = false
					}

					resource "lavinmq_shovel" "test_q2q" {
						name       = "vcr_test_q2q"
						vhost      = "/"
						src_uri    = "%[1]s"
						dest_uri   = "%[1]s"
						src_queue  = lavinmq_queue.source_queue.name
						dest_queue = lavinmq_queue.dest_queue.name
					}

					resource "lavinmq_publish_message" "example_message" {
						vhost       = "/"
						exchange    = "amq.default"
						routing_key = lavinmq_queue.source_queue.name
						payload     = "{\"message\": \"VCR test q2q\"}"
						properties = {
							content_type = "application/json"
						}

						depends_on = [
							lavinmq_shovel.test_q2q
						]
					}

					data "lavinmq_queues" "all_queues" {
						vhost = "/"

						depends_on = [
							lavinmq_publish_message.example_message
						]
					}`, testSrcDestURI),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(shovelResourceName, "name", "vcr_test_q2q"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_queue", "test_q2q_source_queue"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_queue", "test_q2q_dest_queue"),
					resource.TestCheckNoResourceAttr(shovelResourceName, "src_exchange"),
					resource.TestCheckNoResourceAttr(shovelResourceName, "dest_exchange"),
					resource.TestCheckResourceAttrSet(queueDataSourceName, "queues.#"),
					resource.TestCheckTypeSetElemNestedAttrs(queueDataSourceName, "queues.*", map[string]string{
						"name":  "test_q2q_dest_queue",
						"ready": "1",
					}),
				),
			},
		},
	})
}

func TestAccShovel_QueueToExchange(t *testing.T) {
	t.Parallel()
	shovelResourceName := "lavinmq_shovel.test_q2e"
	exchangeDataSourceName := "data.lavinmq_exchanges.dest_exchange"

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
					resource "lavinmq_queue" "source_queue" {
						name        = "test_q2e_source_queue"
						vhost       = "/"
						durable     = true
						auto_delete = false
					}

					resource "lavinmq_exchange" "dest_exchange" {
						name        = "test_q2e_dest_exchange"
						vhost       = "/"
						type        = "direct"
						auto_delete = false
						durable     = true
					}

					resource "lavinmq_shovel" "test_q2e" {
						name             = "vcr_test_q2e"
						vhost            = "/"
						src_uri          = "%[1]s"
						dest_uri         = "%[1]s"
						src_queue        = lavinmq_queue.source_queue.name
						dest_exchange    = lavinmq_exchange.dest_exchange.name
						dest_exchange_key = "routing.key"
					}

					data "lavinmq_exchanges" "dest_exchange" {
						vhost = "/"

						depends_on = [
							lavinmq_shovel.test_q2e
						]
					}`, testSrcDestURI),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(shovelResourceName, "name", "vcr_test_q2e"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_queue", "test_q2e_source_queue"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_exchange", "test_q2e_dest_exchange"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_exchange_key", "routing.key"),
					resource.TestCheckNoResourceAttr(shovelResourceName, "src_exchange"),
					resource.TestCheckNoResourceAttr(shovelResourceName, "dest_queue"),
					resource.TestCheckResourceAttrSet(exchangeDataSourceName, "exchanges.#"),
					resource.TestCheckTypeSetElemNestedAttrs(exchangeDataSourceName, "exchanges.*", map[string]string{
						"name":                      "test_q2e_dest_exchange",
						"message_stats.publish_in":  "0",
						"message_stats.publish_out": "0",
						"message_stats.unroutable":  "0",
					}),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_queue" "source_queue" {
						name        = "test_q2e_source_queue"
						vhost       = "/"
						durable     = true
						auto_delete = false
					}

					resource "lavinmq_exchange" "dest_exchange" {
						name        = "test_q2e_dest_exchange"
						vhost       = "/"
						type        = "direct"
						auto_delete = false
						durable     = true
					}

					resource "lavinmq_shovel" "test_q2e" {
						name             = "vcr_test_q2e"
						vhost            = "/"
						src_uri          = "%[1]s"
						dest_uri         = "%[1]s"
						src_queue        = lavinmq_queue.source_queue.name
						dest_exchange    = lavinmq_exchange.dest_exchange.name
						dest_exchange_key = "routing.key"
					}

					resource "lavinmq_publish_message" "example_message" {
						vhost       = "/"
						exchange    = "amq.default"
						routing_key = lavinmq_queue.source_queue.name
						payload     = "{\"message\": \"VCR test q2e\"}"
						properties = {
							content_type = "application/json"
						}

						depends_on = [
							lavinmq_shovel.test_q2e
						]
					}

					data "lavinmq_exchanges" "dest_exchange" {
						vhost = "/"

						depends_on = [
							lavinmq_publish_message.example_message
						]
					}`, testSrcDestURI),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(shovelResourceName, "name", "vcr_test_q2e"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_queue", "test_q2e_source_queue"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_exchange", "test_q2e_dest_exchange"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_exchange_key", "routing.key"),
					resource.TestCheckNoResourceAttr(shovelResourceName, "src_exchange"),
					resource.TestCheckNoResourceAttr(shovelResourceName, "dest_queue"),
					resource.TestCheckResourceAttrSet(exchangeDataSourceName, "exchanges.#"),
					resource.TestCheckTypeSetElemNestedAttrs(exchangeDataSourceName, "exchanges.*", map[string]string{
						"name":                      "test_q2e_dest_exchange",
						"message_stats.publish_in":  "1",
						"message_stats.publish_out": "0",
						"message_stats.unroutable":  "1",
					}),
				),
			},
		},
	})
}

func TestAccShovel_ExchangeToQueue(t *testing.T) {
	t.Parallel()
	shovelResourceName := "lavinmq_shovel.test_e2q"
	queueDataSourceName := "data.lavinmq_queues.all_queues"

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
					resource "lavinmq_exchange" "source_exchange" {
							name        = "test_e2q_source_exchange"
							vhost       = "/"
							type        = "topic"
							auto_delete = false
							durable     = true
					}

					resource "lavinmq_queue" "dest_queue" {
							name        = "test_e2q_dest_queue"
							vhost       = "/"
							durable     = true
							auto_delete = false
					}

					resource "lavinmq_shovel" "test_e2q" {
							name              = "vcr_test_e2q"
							vhost             = "/"
							src_uri           = "%[1]s"
							dest_uri          = "%[1]s"
							src_exchange      = lavinmq_exchange.source_exchange.name
							src_exchange_key  = "test.#"
							dest_queue        = lavinmq_queue.dest_queue.name
					}

					data "lavinmq_queues" "all_queues" {
							vhost = "/"

							depends_on = [
									lavinmq_shovel.test_e2q
							]
					}`, testSrcDestURI),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(shovelResourceName, "name", "vcr_test_e2q"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_exchange", "test_e2q_source_exchange"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_exchange_key", "test.#"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_queue", "test_e2q_dest_queue"),
					resource.TestCheckNoResourceAttr(shovelResourceName, "src_queue"),
					resource.TestCheckNoResourceAttr(shovelResourceName, "dest_exchange"),
					resource.TestCheckResourceAttrSet(queueDataSourceName, "queues.#"),
					resource.TestCheckTypeSetElemNestedAttrs(queueDataSourceName, "queues.*", map[string]string{
						"name":  "test_e2q_dest_queue",
						"ready": "0",
					}),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_exchange" "source_exchange" {
							name        = "test_e2q_source_exchange"
							vhost       = "/"
							type        = "topic"
							auto_delete = false
							durable     = true
					}

					resource "lavinmq_queue" "dest_queue" {
							name        = "test_e2q_dest_queue"
							vhost       = "/"
							durable     = true
							auto_delete = false
					}

					resource "lavinmq_shovel" "test_e2q" {
							name              = "vcr_test_e2q"
							vhost             = "/"
							src_uri           = "%[1]s"
							dest_uri          = "%[1]s"
							src_exchange      = lavinmq_exchange.source_exchange.name
							src_exchange_key  = "test.#"
							dest_queue        = lavinmq_queue.dest_queue.name
					}

					resource "lavinmq_publish_message" "example_message" {
							vhost       = "/"
							exchange    = lavinmq_exchange.source_exchange.name
							routing_key = "test.message"
							payload     = "{\"message\": \"VCR test e2q\"}"
							properties = {
									content_type = "application/json"
							}

							depends_on = [
									lavinmq_shovel.test_e2q
							]
					}

					data "lavinmq_queues" "all_queues" {
							vhost = "/"

							depends_on = [
									lavinmq_publish_message.example_message
							]
					}`, testSrcDestURI),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(shovelResourceName, "name", "vcr_test_e2q"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_exchange", "test_e2q_source_exchange"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_exchange_key", "test.#"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_queue", "test_e2q_dest_queue"),
					resource.TestCheckNoResourceAttr(shovelResourceName, "src_queue"),
					resource.TestCheckNoResourceAttr(shovelResourceName, "dest_exchange"),
					resource.TestCheckResourceAttrSet(queueDataSourceName, "queues.#"),
					resource.TestCheckTypeSetElemNestedAttrs(queueDataSourceName, "queues.*", map[string]string{
						"name":  "test_e2q_dest_queue",
						"ready": "1",
					}),
				),
			},
		},
	})
}

func TestAccShovel_ExchangeToExchange(t *testing.T) {
	t.Parallel()
	shovelResourceName := "lavinmq_shovel.test_e2e"
	exchangeDataSourceName := "data.lavinmq_exchanges.dest_exchange"

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
					resource "lavinmq_exchange" "source_exchange" {
							name        = "test_e2e_source_exchange"
							vhost       = "/"
							type        = "topic"
							auto_delete = false
							durable     = true
					}

					resource "lavinmq_exchange" "dest_exchange" {
							name        = "test_e2e_dest_exchange"
							vhost       = "/"
							type        = "topic"
							auto_delete = false
							durable     = true
					}

					resource "lavinmq_shovel" "test_e2e" {
							name              = "vcr_test_e2e"
							vhost             = "/"
							src_uri           = "%[1]s"
							dest_uri          = "%[1]s"
							src_exchange      = lavinmq_exchange.source_exchange.name
							src_exchange_key  = "source.#"
							dest_exchange     = lavinmq_exchange.dest_exchange.name
							dest_exchange_key = "dest.key"
					}

					data "lavinmq_exchanges" "dest_exchange" {
							vhost = "/"

							depends_on = [
									lavinmq_shovel.test_e2e
							]
					}`, testSrcDestURI),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(shovelResourceName, "name", "vcr_test_e2e"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_exchange", "test_e2e_source_exchange"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_exchange_key", "source.#"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_exchange", "test_e2e_dest_exchange"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_exchange_key", "dest.key"),
					resource.TestCheckNoResourceAttr(shovelResourceName, "src_queue"),
					resource.TestCheckNoResourceAttr(shovelResourceName, "dest_queue"),
					resource.TestCheckResourceAttrSet(exchangeDataSourceName, "exchanges.#"),
					resource.TestCheckTypeSetElemNestedAttrs(exchangeDataSourceName, "exchanges.*", map[string]string{
						"name":                      "test_e2e_dest_exchange",
						"message_stats.publish_in":  "0",
						"message_stats.publish_out": "0",
					}),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_exchange" "source_exchange" {
							name        = "test_e2e_source_exchange"
							vhost       = "/"
							type        = "topic"
							auto_delete = false
							durable     = true
					}

					resource "lavinmq_exchange" "dest_exchange" {
							name        = "test_e2e_dest_exchange"
							vhost       = "/"
							type        = "topic"
							auto_delete = false
							durable     = true
					}

					resource "lavinmq_shovel" "test_e2e" {
							name              = "vcr_test_e2e"
							vhost             = "/"
							src_uri           = "%[1]s"
							dest_uri          = "%[1]s"
							src_exchange      = lavinmq_exchange.source_exchange.name
							src_exchange_key  = "source.#"
							dest_exchange     = lavinmq_exchange.dest_exchange.name
							dest_exchange_key = "dest.key"
					}

					resource "lavinmq_publish_message" "example_message" {
							vhost       = "/"
							exchange    = lavinmq_exchange.source_exchange.name
							routing_key = "source.message"
							payload     = "{\"message\": \"VCR test e2e\"}"
							properties = {
									content_type = "application/json"
							}

							depends_on = [
									lavinmq_shovel.test_e2e
							]
					}

					data "lavinmq_exchanges" "dest_exchange" {
							vhost = "/"

							depends_on = [
									lavinmq_publish_message.example_message
							]
					}`, testSrcDestURI),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(shovelResourceName, "name", "vcr_test_e2e"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_exchange", "test_e2e_source_exchange"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_exchange_key", "source.#"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_exchange", "test_e2e_dest_exchange"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_exchange_key", "dest.key"),
					resource.TestCheckNoResourceAttr(shovelResourceName, "src_queue"),
					resource.TestCheckNoResourceAttr(shovelResourceName, "dest_queue"),
					resource.TestCheckResourceAttrSet(exchangeDataSourceName, "exchanges.#"),
					resource.TestCheckTypeSetElemNestedAttrs(exchangeDataSourceName, "exchanges.*", map[string]string{
						"name":                      "test_e2e_dest_exchange",
						"message_stats.publish_in":  "1",
						"message_stats.publish_out": "0",
					}),
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
