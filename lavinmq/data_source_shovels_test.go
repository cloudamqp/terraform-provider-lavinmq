package lavinmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceShovels_Basic(t *testing.T) {
	t.Parallel()
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "lavinmq_shovel" "test1" {
  name       = "vcr_test_shovel_ds_1"
  vhost      = "/"
  src_uri    = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri   = "amqp://guest:guest@localhost:5672/%2f"
  src_queue  = "source_queue"
  dest_queue = "dest_queue"
}

resource "lavinmq_shovel" "test2" {
  name              = "vcr_test_shovel_ds_2"
  vhost             = "/"
  src_uri           = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri          = "amqp://guest:guest@localhost:5672/%2f"
  src_exchange      = "source_exchange"
  src_exchange_key  = "source.key"
  dest_queue        = "dest_queue"
}

data "lavinmq_shovels" "all" {
  depends_on = [lavinmq_shovel.test1, lavinmq_shovel.test2]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_shovels.all", "shovels.*", map[string]string{
						"name":       "vcr_test_shovel_ds_1",
						"vhost":      "/",
						"src_queue":  "source_queue",
						"dest_queue": "dest_queue",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_shovels.all", "shovels.*", map[string]string{
						"name":             "vcr_test_shovel_ds_2",
						"vhost":            "/",
						"src_exchange":     "source_exchange",
						"src_exchange_key": "source.key",
						"dest_queue":       "dest_queue",
					}),
				),
			},
		},
	})
}
