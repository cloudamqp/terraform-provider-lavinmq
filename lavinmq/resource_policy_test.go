package lavinmq

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPolicy_Import(t *testing.T) {
	t.Parallel()
	policyResourceName := "lavinmq_policy.test_policy"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_policy" "test_policy" {
            name     = "vcr_test_policy_import"
            vhost    = "/"
            pattern  = "^vcr_test"
            definition = {
            }
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyResourceName, "name", "vcr_test_policy_import"),
					resource.TestCheckResourceAttr(policyResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(policyResourceName, "pattern", "^vcr_test"),
					resource.TestCheckResourceAttr(policyResourceName, "priority", "0"),
					resource.TestCheckResourceAttr(policyResourceName, "apply_to", "all"),
				),
			},
			{
				ResourceName:                         policyResourceName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateId:                        "/@vcr_test_policy_import",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

func TestAccPolicy_Update(t *testing.T) {
	t.Parallel()
	policyResourceName := "lavinmq_policy.test_policy"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_policy" "test_policy" {
            name     = "vcr_test_policy_update"
            vhost    = "/"
            pattern  = "^vcr_test"
            definition = {
              "message-ttl" = 60000
              "max-length"  = 500
            }
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyResourceName, "name", "vcr_test_policy_update"),
					resource.TestCheckResourceAttr(policyResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(policyResourceName, "pattern", "^vcr_test"),
					resource.TestCheckResourceAttr(policyResourceName, "definition.message-ttl", "60000"),
					resource.TestCheckResourceAttr(policyResourceName, "definition.max-length", "500"),
					resource.TestCheckResourceAttr(policyResourceName, "priority", "0"),
					resource.TestCheckResourceAttr(policyResourceName, "apply_to", "all"),
				),
			},
			{
				Config: `
          resource "lavinmq_policy" "test_policy" {
            name     = "vcr_test_policy_update"
            vhost    = "/"
            pattern  = "^vcr_test"
            definition = {
              "message-ttl" = 3600000
              "max-length"  = 1000
            }
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyResourceName, "name", "vcr_test_policy_update"),
					resource.TestCheckResourceAttr(policyResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(policyResourceName, "pattern", "^vcr_test"),
					resource.TestCheckResourceAttr(policyResourceName, "definition.message-ttl", "3600000"),
					resource.TestCheckResourceAttr(policyResourceName, "definition.max-length", "1000"),
					resource.TestCheckResourceAttr(policyResourceName, "priority", "0"),
					resource.TestCheckResourceAttr(policyResourceName, "apply_to", "all"),
				),
			},
		},
	})
}

func TestAccPolicy_AddDefinitions(t *testing.T) {
	t.Parallel()
	policyResourceName := "lavinmq_policy.dead_letter_policy"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_policy" "dead_letter_policy" {
            name     = "vcr_test_dl_policy"
            vhost    = "/"
            pattern  = "^dl_test"
            definition = {
            }
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyResourceName, "name", "vcr_test_dl_policy"),
					resource.TestCheckResourceAttr(policyResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(policyResourceName, "pattern", "^dl_test"),
					resource.TestCheckResourceAttr(policyResourceName, "apply_to", "all"),
				),
			},
			{
				Config: `
          resource "lavinmq_policy" "dead_letter_policy" {
            name     = "vcr_test_dl_policy"
            vhost    = "/"
            pattern  = "^dl_test"
            definition = {
              "dead-letter-exchange"    = "updated_dlx_exchange"
              "dead-letter-routing-key" = "updated_dlx_routing_key"
            }
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyResourceName, "name", "vcr_test_dl_policy"),
					resource.TestCheckResourceAttr(policyResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(policyResourceName, "pattern", "^dl_test"),
					resource.TestCheckResourceAttr(policyResourceName, "definition.dead-letter-exchange", "updated_dlx_exchange"),
					resource.TestCheckResourceAttr(policyResourceName, "definition.dead-letter-routing-key", "updated_dlx_routing_key"),
					resource.TestCheckResourceAttr(policyResourceName, "apply_to", "all"),
				),
			},
		},
	})
}

func TestAccPolicy_InvalidApplyTo(t *testing.T) {
	t.Parallel()
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_policy" "test_policy" {
            name     = "vcr_test_policy_invalid"
            vhost    = "/"
            pattern  = "^vcr_test"
            apply_to = "invalid_value"
            definition = {
            }
          }`,
				ExpectError: regexp.MustCompile(`Attribute apply_to value must be one of:.*"all".*"exchanges".*"queues"`),
			},
		},
	})
}
