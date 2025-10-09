package lavinmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceUsers_Basic(t *testing.T) {
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_user" "test1" {
            name = "terraform-lavinmq-user-test-1"
            password = "test-password-1"
            tags = ["monitoring"]
          }

          resource "lavinmq_user" "test2" {
            name = "terraform-lavinmq-user-test-2"
            password = "test-password-2"
            tags = ["management", "policymaker"]
          }

          data "lavinmq_users" "all" {
            depends_on = [lavinmq_user.test1, lavinmq_user.test2]
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_users.all", "users.*", map[string]string{
						"name":   "terraform-lavinmq-user-test-1",
						"tags.#": "1",
						"tags.0": "monitoring",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_users.all", "users.*", map[string]string{
						"name":   "terraform-lavinmq-user-test-2",
						"tags.#": "2",
						"tags.0": "management",
						"tags.1": "policymaker",
					}),
				),
			},
		},
	})
}

func TestAccDataSourceUsers_DefaultUser(t *testing.T) {
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          data "lavinmq_users" "all" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.lavinmq_users.all", "users.#"),
				),
			},
		},
	})
}
