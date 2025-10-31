package lavinmq

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFederationUpstream_Import(t *testing.T) {
	t.Parallel()
	resourceName := "lavinmq_federation_upstream.test"

	// Set sanitized value for playback and use real value for recording
	testUpstreamURI := "TEST_AMQP_URI"
	if os.Getenv("LAVINMQ_RECORD") != "" {
		testUpstreamURI = os.Getenv("TEST_AMQP_URI")
	}

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_federation_upstream" "test" {
							name     = "vcr_test_federation_upstream_import"
							vhost    = "/"
							uri      = "%[1]s"
							exchange = "upstream-exchange"
					}`, testUpstreamURI),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vcr_test_federation_upstream_import"),
					resource.TestCheckResourceAttr(resourceName, "vhost", "/"),
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

	// Set sanitized value for playback and use real value for recording
	testUpstreamURI := "TEST_AMQP_URI"
	if os.Getenv("LAVINMQ_RECORD") != "" {
		testUpstreamURI = os.Getenv("TEST_AMQP_URI")
	}

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_federation_upstream" "test" {
							name           = "vcr_test_federation_upstream_update"
							vhost          = "/"
							uri            = "%[1]s"
							exchange       = "upstream-exchange"
							prefetch_count = 500
							ack_mode       = "on-publish"
					}`, testUpstreamURI),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vcr_test_federation_upstream_update"),
					resource.TestCheckResourceAttr(resourceName, "prefetch_count", "500"),
					resource.TestCheckResourceAttr(resourceName, "ack_mode", "on-publish"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_federation_upstream" "test" {
							name            = "vcr_test_federation_upstream_update"
							vhost           = "/"
							uri             = "%[1]s"
							exchange        = "upstream-exchange"
							prefetch_count  = 2000
							ack_mode        = "no-ack"
							reconnect_delay = 10
							max_hops        = 3
					}`, testUpstreamURI),
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
	exchangeDataSourceName := "data.lavinmq_exchanges.all_exchanges"

	// Set sanitized value for playback and use real value for recording
	testUpstreamURI := "TEST_AMQP_URI"
	if os.Getenv("LAVINMQ_RECORD") != "" {
		testUpstreamURI = os.Getenv("TEST_AMQP_URI")
	}

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_exchange" "local_exchange" {
							name        = "test_fed_local_exchange"
							vhost       = "/"
							type        = "topic"
							auto_delete = false
							durable     = true
					}

					resource "lavinmq_federation_upstream" "test_exchange" {
							name     = "vcr_test_exchange_federation"
							vhost    = "/"
							uri      = "%[1]s"
							exchange = "federated-exchange"
							max_hops = 2
					}

					resource "lavinmq_policy" "federation_policy" {
							name       = "test_fed_policy"
							vhost      = "/"
							pattern    = "^test_fed_.*"
							apply_to   = "exchanges"
							priority   = 10
							definition = {
									"federation-upstream" = lavinmq_federation_upstream.test_exchange.name
							}

							depends_on = [
									lavinmq_exchange.local_exchange,
									lavinmq_federation_upstream.test_exchange
							]
					}

					data "lavinmq_exchanges" "all_exchanges" {
							vhost = "/"

							depends_on = [
									lavinmq_policy.federation_policy
							]
					}`, testUpstreamURI),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vcr_test_exchange_federation"),
					resource.TestCheckResourceAttr(resourceName, "exchange", "federated-exchange"),
					resource.TestCheckResourceAttr(resourceName, "max_hops", "2"),
					resource.TestCheckNoResourceAttr(resourceName, "queue"),
					resource.TestCheckResourceAttrSet(exchangeDataSourceName, "exchanges.#"),
					resource.TestCheckTypeSetElemNestedAttrs(exchangeDataSourceName, "exchanges.*", map[string]string{
						"name":                      "test_fed_local_exchange",
						"message_stats.publish_in":  "0",
						"message_stats.publish_out": "0",
					}),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_exchange" "local_exchange" {
							name        = "test_fed_local_exchange"
							vhost       = "/"
							type        = "topic"
							auto_delete = false
							durable     = true
					}

					resource "lavinmq_federation_upstream" "test_exchange" {
							name     = "vcr_test_exchange_federation"
							vhost    = "/"
							uri      = "%[1]s"
							exchange = "federated-exchange"
							max_hops = 2
					}

					resource "lavinmq_policy" "federation_policy" {
							name       = "test_fed_policy"
							vhost      = "/"
							pattern    = "^test_fed_.*"
							apply_to   = "exchanges"
							priority   = 10
							definition = {
									"federation-upstream" = lavinmq_federation_upstream.test_exchange.name
							}

							depends_on = [
									lavinmq_exchange.local_exchange,
									lavinmq_federation_upstream.test_exchange
							]
					}

					resource "lavinmq_publish_message" "test_message" {
							vhost       = "/"
							exchange    = lavinmq_exchange.local_exchange.name
							routing_key = "test.federation"
							payload     = "{\"message\": \"VCR test exchange federation\"}"
							properties = {
									content_type = "application/json"
							}

							depends_on = [
									lavinmq_policy.federation_policy
							]
					}

					data "lavinmq_exchanges" "all_exchanges" {
							vhost = "/"

							depends_on = [
									lavinmq_publish_message.test_message
							]
					}`, testUpstreamURI),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vcr_test_exchange_federation"),
					resource.TestCheckResourceAttr(resourceName, "exchange", "federated-exchange"),
					resource.TestCheckResourceAttr(resourceName, "max_hops", "2"),
					resource.TestCheckNoResourceAttr(resourceName, "queue"),
					resource.TestCheckResourceAttrSet(exchangeDataSourceName, "exchanges.#"),
					resource.TestCheckTypeSetElemNestedAttrs(exchangeDataSourceName, "exchanges.*", map[string]string{
						"name":                     "test_fed_local_exchange",
						"message_stats.publish_in": "1",
					}),
				),
			},
		},
	})
}

func TestAccFederationUpstream_QueueFederation(t *testing.T) {
	t.Parallel()
	resourceName := "lavinmq_federation_upstream.test_queue"
	queueDataSourceName := "data.lavinmq_queues.all_queues"

	// Set sanitized value for playback and use real value for recording
	testUpstreamURI := "TEST_AMQP_URI"
	if os.Getenv("LAVINMQ_RECORD") != "" {
		testUpstreamURI = os.Getenv("TEST_AMQP_URI")
	}

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_queue" "local_queue" {
							name        = "test_fed_local_queue"
							vhost       = "/"
							durable     = true
							auto_delete = false
					}

					resource "lavinmq_federation_upstream" "test_queue" {
							name  = "vcr_test_queue_federation"
							vhost = "/"
							uri   = "%[1]s"
							queue = "federated-queue"
					}

					resource "lavinmq_policy" "federation_policy" {
							name       = "test_fed_queue_policy"
							vhost      = "/"
							pattern    = "^test_fed_.*"
							apply_to   = "queues"
							priority   = 10
							definition = {
									"federation-upstream" = lavinmq_federation_upstream.test_queue.name
							}

							depends_on = [
									lavinmq_queue.local_queue,
									lavinmq_federation_upstream.test_queue
							]
					}

					data "lavinmq_queues" "all_queues" {
							vhost = "/"

							depends_on = [
									lavinmq_policy.federation_policy
							]
					}`, testUpstreamURI),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vcr_test_queue_federation"),
					resource.TestCheckResourceAttr(resourceName, "queue", "federated-queue"),
					resource.TestCheckNoResourceAttr(resourceName, "exchange"),
					resource.TestCheckResourceAttrSet(queueDataSourceName, "queues.#"),
					resource.TestCheckTypeSetElemNestedAttrs(queueDataSourceName, "queues.*", map[string]string{
						"name":  "test_fed_local_queue",
						"ready": "0",
					}),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_queue" "local_queue" {
							name        = "test_fed_local_queue"
							vhost       = "/"
							durable     = true
							auto_delete = false
					}

					resource "lavinmq_federation_upstream" "test_queue" {
							name  = "vcr_test_queue_federation"
							vhost = "/"
							uri   = "%[1]s"
							queue = "federated-queue"
					}

					resource "lavinmq_policy" "federation_policy" {
							name       = "test_fed_queue_policy"
							vhost      = "/"
							pattern    = "^test_fed_.*"
							apply_to   = "queues"
							priority   = 10
							definition = {
									"federation-upstream" = lavinmq_federation_upstream.test_queue.name
							}

							depends_on = [
									lavinmq_queue.local_queue,
									lavinmq_federation_upstream.test_queue
							]
					}

					resource "lavinmq_publish_message" "test_message" {
							vhost       = "/"
							exchange    = "amq.default"
							routing_key = lavinmq_queue.local_queue.name
							payload     = "{\"message\": \"VCR test queue federation\"}"
							properties = {
									content_type = "application/json"
							}

							depends_on = [
									lavinmq_policy.federation_policy
							]
					}

					data "lavinmq_queues" "all_queues" {
							vhost = "/"

							depends_on = [
									lavinmq_publish_message.test_message
							]
					}`, testUpstreamURI),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vcr_test_queue_federation"),
					resource.TestCheckResourceAttr(resourceName, "queue", "federated-queue"),
					resource.TestCheckNoResourceAttr(resourceName, "exchange"),
					resource.TestCheckResourceAttrSet(queueDataSourceName, "queues.#"),
					resource.TestCheckTypeSetElemNestedAttrs(queueDataSourceName, "queues.*", map[string]string{
						"name":  "test_fed_local_queue",
						"ready": "1",
					}),
				),
			},
		},
	})
}

func TestAccFederationUpstream_WithTTL(t *testing.T) {
	t.Parallel()
	resourceName := "lavinmq_federation_upstream.test_ttl"

	// Set sanitized value for playback and use real value for recording
	testUpstreamURI := "TEST_AMQP_URI"
	if os.Getenv("LAVINMQ_RECORD") != "" {
		testUpstreamURI = os.Getenv("TEST_AMQP_URI")
	}

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "lavinmq_federation_upstream" "test_ttl" {
							name        = "vcr_test_federation_ttl"
							vhost       = "/"
							uri         = "%[1]s"
							exchange    = "ttl-exchange"
							expires     = 3600000
							message_ttl = 60000
					}`, testUpstreamURI),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vcr_test_federation_ttl"),
					resource.TestCheckResourceAttr(resourceName, "expires", "3600000"),
					resource.TestCheckResourceAttr(resourceName, "message_ttl", "60000"),
				),
			},
		},
	})
}
