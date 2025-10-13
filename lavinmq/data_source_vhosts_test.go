package lavinmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceVhosts_Basic(t *testing.T) {
	t.Parallel()
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_vhost" "test1" {
            name = "terraform-lavinmq-test-1"
          }

          resource "lavinmq_vhost" "test2" {
            name = "terraform-lavinmq-test-2"
          }

          data "lavinmq_vhosts" "all" {
            depends_on = [lavinmq_vhost.test1, lavinmq_vhost.test2]
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_vhosts.all", "vhosts.*", map[string]string{
						"name": "terraform-lavinmq-test-1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_vhosts.all", "vhosts.*", map[string]string{
						"name": "terraform-lavinmq-test-2",
					}),
				),
			},
		},
	})
}

func TestAccDataSourceVhosts_DefaultVhost(t *testing.T) {
	t.Parallel()
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          data "lavinmq_vhosts" "all" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.lavinmq_vhosts.all", "vhosts.#"),
				),
			},
		},
	})
}
