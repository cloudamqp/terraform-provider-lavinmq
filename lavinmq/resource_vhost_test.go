package lavinmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVhost_Basic(t *testing.T) {
	vhostResourceName := "lavinmq_vhost.vcr_test"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_vhost" "vcr_test" {
            name = "vcr_test"
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(vhostResourceName, "name", "vcr_test"),
					resource.TestCheckNoResourceAttr(vhostResourceName, "max_connections"),
					resource.TestCheckNoResourceAttr(vhostResourceName, "max_queues"),
				),
			},
			{
				ResourceName:                         vhostResourceName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateId:                        "vcr_test",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
			{
				Config: `
          resource "lavinmq_vhost" "vcr_test" {
            name            = "vcr_test"
            max_connections = 100
            max_queues      = 30
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(vhostResourceName, "name", "vcr_test"),
					resource.TestCheckResourceAttr(vhostResourceName, "max_connections", "100"),
					resource.TestCheckResourceAttr(vhostResourceName, "max_queues", "30"),
				),
			},
			{
				Config: `
          resource "lavinmq_vhost" "vcr_test" {
            name            = "vcr_test"
            max_connections = 100
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(vhostResourceName, "name", "vcr_test"),
					resource.TestCheckResourceAttr(vhostResourceName, "max_connections", "100"),
					resource.TestCheckNoResourceAttr(vhostResourceName, "max_queues"),
				),
			},
			{
				Config: `
          resource "lavinmq_vhost" "vcr_test" {
            name       = "vcr_test"
            max_queues = 30
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(vhostResourceName, "name", "vcr_test"),
					resource.TestCheckResourceAttr(vhostResourceName, "max_queues", "30"),
					resource.TestCheckNoResourceAttr(vhostResourceName, "max_connections"),
				),
			},
		},
	})
}
