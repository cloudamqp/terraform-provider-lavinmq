package lavinmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourcePolicies_Basic(t *testing.T) {
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_policy" "test1" {
            name     = "terraform-lavinmq-test-policy-1"
            vhost    = "/"
            pattern  = "^test1"
            priority = 0
            apply_to = "queues"
            definition = {
              "message-ttl" = 60000
            }
          }

          resource "lavinmq_policy" "test2" {
            name     = "terraform-lavinmq-test-policy-2"
            vhost    = "/"
            pattern  = "^test2"
            priority = 1
            apply_to = "exchanges"
            definition = {
              "max-length" = 1000
            }
          }

          data "lavinmq_policies" "all" {
            depends_on = [lavinmq_policy.test1, lavinmq_policy.test2]
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_policies.all", "policies.*", map[string]string{
						"name":     "terraform-lavinmq-test-policy-1",
						"vhost":    "/",
						"pattern":  "^test1",
						"priority": "0",
						"apply_to": "queues",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_policies.all", "policies.*", map[string]string{
						"name":     "terraform-lavinmq-test-policy-2",
						"vhost":    "/",
						"pattern":  "^test2",
						"priority": "1",
						"apply_to": "exchanges",
					}),
				),
			},
		},
	})
}
