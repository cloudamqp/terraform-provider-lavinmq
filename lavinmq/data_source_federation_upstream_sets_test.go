package lavinmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceFederationUpstreamSets_Basic(t *testing.T) {
	t.Parallel()
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "lavinmq_federation_upstream" "upstream1" {
  name     = "vcr_test_ds_upstream-1"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstream1:5672/%2f"
  exchange = "exchange1"
}

resource "lavinmq_federation_upstream" "upstream2" {
  name     = "vcr_test_ds_upstream-2"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstream2:5672/%2f"
  exchange = "exchange2"
}

resource "lavinmq_federation_upstream" "upstreama" {
  name     = "vcr_test_ds_upstream-a"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstreama:5672/%2f"
  exchange = "exchangea"
}

resource "lavinmq_federation_upstream" "upstreamb" {
  name     = "vcr_test_ds_upstream-b"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstreamb:5672/%2f"
  exchange = "exchangeb"
}

resource "lavinmq_federation_upstream" "upstreamc" {
  name     = "vcr_test_ds_upstream-c"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstreamc:5672/%2f"
  exchange = "exchangec"
}

resource "lavinmq_federation_upstream_set" "test1" {
  name   = "vcr_test_federation_set_ds_1"
  vhost  = "/"
  upstreams = [
    lavinmq_federation_upstream.upstream1.name,
    lavinmq_federation_upstream.upstream2.name
  ]
}

resource "lavinmq_federation_upstream_set" "test2" {
  name   = "vcr_test_federation_set_ds_2"
  vhost  = "/"
  upstreams = [
    lavinmq_federation_upstream.upstreama.name,
    lavinmq_federation_upstream.upstreamb.name,
    lavinmq_federation_upstream.upstreamc.name
  ]
}

data "lavinmq_federation_upstream_sets" "all" {
  depends_on = [lavinmq_federation_upstream_set.test1, lavinmq_federation_upstream_set.test2]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_federation_upstream_sets.all", "federation_upstream_sets.*", map[string]string{
						"name":        "vcr_test_federation_set_ds_1",
						"vhost":       "/",
						"upstreams.#": "2",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_federation_upstream_sets.all", "federation_upstream_sets.*", map[string]string{
						"name":        "vcr_test_federation_set_ds_2",
						"vhost":       "/",
						"upstreams.#": "3",
					}),
				),
			},
		},
	})
}
