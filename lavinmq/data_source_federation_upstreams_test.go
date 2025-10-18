package lavinmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceFederationUpstreams_Basic(t *testing.T) {
	t.Parallel()
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "lavinmq_federation_upstream" "test1" {
  name     = "vcr_test_federation_ds_1"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstream1:5672/%2f"
  exchange = "exchange1"
}

resource "lavinmq_federation_upstream" "test2" {
  name           = "vcr_test_federation_ds_2"
  vhost          = "/"
  uri            = "amqp://guest:guest@upstream2:5672/%2f"
  queue          = "queue1"
  prefetch_count = 500
}

data "lavinmq_federation_upstreams" "all" {
  depends_on = [lavinmq_federation_upstream.test1, lavinmq_federation_upstream.test2]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_federation_upstreams.all", "federation_upstreams.*", map[string]string{
						"name":     "vcr_test_federation_ds_1",
						"vhost":    "/",
						"exchange": "exchange1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_federation_upstreams.all", "federation_upstreams.*", map[string]string{
						"name":           "vcr_test_federation_ds_2",
						"vhost":          "/",
						"queue":          "queue1",
						"prefetch_count": "500",
					}),
				),
			},
		},
	})
}
