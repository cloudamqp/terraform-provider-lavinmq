package lavinmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFederationUpstreamSet_Import(t *testing.T) {
	t.Parallel()
	resourceName := "lavinmq_federation_upstream_set.test"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "lavinmq_federation_upstream" "upstream1" {
  name     = "vcr_test_upstream-1"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstream1:5672/%2f"
  exchange = "exchange1"
}

resource "lavinmq_federation_upstream" "upstream2" {
  name     = "vcr_test_upstream-2"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstream2:5672/%2f"
  exchange = "exchange2"
}

resource "lavinmq_federation_upstream_set" "test" {
  name   = "vcr_test_upstream_set_import"
  vhost  = "/"
  upstreams = [
    lavinmq_federation_upstream.upstream1.name,
    lavinmq_federation_upstream.upstream2.name
  ]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vcr_test_upstream_set_import"),
					resource.TestCheckResourceAttr(resourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(resourceName, "upstreams.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "upstreams.0", "vcr_test_upstream-1"),
					resource.TestCheckResourceAttr(resourceName, "upstreams.1", "vcr_test_upstream-2"),
				),
			},
			{
				ResourceName:                         resourceName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateId:                        "/@vcr_test_upstream_set_import",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

func TestAccFederationUpstreamSet_Update(t *testing.T) {
	t.Parallel()
	resourceName := "lavinmq_federation_upstream_set.test"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "lavinmq_federation_upstream" "upstream1" {
  name     = "vcr_test_update_upstream-1"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstream1:5672/%2f"
  exchange = "exchange1"
}

resource "lavinmq_federation_upstream" "upstream2" {
  name     = "vcr_test_update_upstream-2"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstream2:5672/%2f"
  exchange = "exchange2"
}

resource "lavinmq_federation_upstream_set" "test" {
  name   = "vcr_test_upstream_set_update"
  vhost  = "/"
  upstreams = [
    lavinmq_federation_upstream.upstream1.name,
    lavinmq_federation_upstream.upstream2.name
  ]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vcr_test_upstream_set_update"),
					resource.TestCheckResourceAttr(resourceName, "upstreams.#", "2"),
				),
			},
			{
				Config: `
resource "lavinmq_federation_upstream" "upstream1" {
  name     = "vcr_test_update_upstream-1"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstream1:5672/%2f"
  exchange = "exchange1"
}

resource "lavinmq_federation_upstream" "upstream2" {
  name     = "vcr_test_update_upstream-2"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstream2:5672/%2f"
  exchange = "exchange2"
}

resource "lavinmq_federation_upstream" "upstream3" {
  name     = "vcr_test_update_upstream-3"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstream3:5672/%2f"
  exchange = "exchange3"
}

resource "lavinmq_federation_upstream_set" "test" {
  name   = "vcr_test_upstream_set_update"
  vhost  = "/"
  upstreams = [
    lavinmq_federation_upstream.upstream1.name,
    lavinmq_federation_upstream.upstream2.name,
    lavinmq_federation_upstream.upstream3.name
  ]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vcr_test_upstream_set_update"),
					resource.TestCheckResourceAttr(resourceName, "upstreams.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "upstreams.2", "vcr_test_update_upstream-3"),
				),
			},
		},
	})
}

func TestAccFederationUpstreamSet_Basic(t *testing.T) {
	t.Parallel()
	resourceName := "lavinmq_federation_upstream_set.test_ha"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "lavinmq_federation_upstream" "us_east" {
  name     = "vcr_test_upstream-us-east"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstream1:5672/%2f"
  exchange = "exchange1"
}

resource "lavinmq_federation_upstream" "us_west" {
  name     = "vcr_test_upstream-us-west"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstream2:5672/%2f"
  exchange = "exchange2"
}

resource "lavinmq_federation_upstream" "eu" {
  name     = "vcr_test_upstream-eu"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstream3:5672/%2f"
  exchange = "exchange3"
}

resource "lavinmq_federation_upstream_set" "test_ha" {
  name   = "vcr_test_ha_set"
  vhost  = "/"
  upstreams = [
    lavinmq_federation_upstream.us_east.name,
    lavinmq_federation_upstream.us_west.name,
    lavinmq_federation_upstream.eu.name
  ]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "vcr_test_ha_set"),
					resource.TestCheckResourceAttr(resourceName, "upstreams.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "upstreams.0", "vcr_test_upstream-us-east"),
					resource.TestCheckResourceAttr(resourceName, "upstreams.1", "vcr_test_upstream-us-west"),
					resource.TestCheckResourceAttr(resourceName, "upstreams.2", "vcr_test_upstream-eu"),
				),
			},
		},
	})
}
