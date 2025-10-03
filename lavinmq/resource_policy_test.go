package lavinmq

import (
	"testing"

	"github.com/cloudamqp/terraform-provider-lavinmq/lavinmq/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPolicy_Import(t *testing.T) {
	var (
		fileNames          = []string{"policies/policy"}
		policyResourceName = "lavinmq_policy.test_policy"

		params = map[string]string{
			"ResourceName":  "test_policy",
			"PolicyName":    "vcr_test_policy_import",
			"PolicyVhost":   "/",
			"PolicyPattern": "^vcr_test",
		}
	)

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyResourceName, "name", params["PolicyName"]),
					resource.TestCheckResourceAttr(policyResourceName, "vhost", params["PolicyVhost"]),
					resource.TestCheckResourceAttr(policyResourceName, "pattern", params["PolicyPattern"]),
					resource.TestCheckResourceAttr(policyResourceName, "priority", "0"),
					resource.TestCheckResourceAttr(policyResourceName, "apply_to", "all"),
				),
			},
			{
				ResourceName:                         policyResourceName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateId:                        params["PolicyVhost"] + "@" + params["PolicyName"],
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

func TestAccPolicy_Update(t *testing.T) {
	var (
		fileNames          = []string{"policies/policy"}
		policyResourceName = "lavinmq_policy.test_policy"

		params = map[string]string{
			"ResourceName":    "test_policy",
			"PolicyName":      "vcr_test_policy_update",
			"PolicyVhost":     "/",
			"PolicyPattern":   "^vcr_test",
			"PolicyTTL":       "60000",
			"PolicyMaxLength": "500",
		}

		fileNamesUpdated = []string{"policies/policy"}
		paramsUpdated    = map[string]string{
			"ResourceName":    "test_policy",
			"PolicyName":      "vcr_test_policy_update",
			"PolicyVhost":     "/",
			"PolicyPattern":   "^vcr_test",
			"PolicyTTL":       "3600000",
			"PolicyMaxLength": "1000",
		}
	)

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyResourceName, "name", params["PolicyName"]),
					resource.TestCheckResourceAttr(policyResourceName, "vhost", params["PolicyVhost"]),
					resource.TestCheckResourceAttr(policyResourceName, "pattern", params["PolicyPattern"]),
					resource.TestCheckResourceAttr(policyResourceName, "definition.message-ttl", params["PolicyTTL"]),
					resource.TestCheckResourceAttr(policyResourceName, "definition.max-length", params["PolicyMaxLength"]),
					resource.TestCheckResourceAttr(policyResourceName, "priority", "0"),
					resource.TestCheckResourceAttr(policyResourceName, "apply_to", "all"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesUpdated, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyResourceName, "name", paramsUpdated["PolicyName"]),
					resource.TestCheckResourceAttr(policyResourceName, "vhost", paramsUpdated["PolicyVhost"]),
					resource.TestCheckResourceAttr(policyResourceName, "pattern", paramsUpdated["PolicyPattern"]),
					resource.TestCheckResourceAttr(policyResourceName, "definition.message-ttl", paramsUpdated["PolicyTTL"]),
					resource.TestCheckResourceAttr(policyResourceName, "definition.max-length", paramsUpdated["PolicyMaxLength"]),
					resource.TestCheckResourceAttr(policyResourceName, "priority", "0"),
					resource.TestCheckResourceAttr(policyResourceName, "apply_to", "all"),
				),
			},
		},
	})
}

func TestAccPolicy_AddDefinitions(t *testing.T) {
	var (
		fileNames          = []string{"policies/policy"}
		policyResourceName = "lavinmq_policy.dead_letter_policy"

		params = map[string]string{
			"ResourceName":  "dead_letter_policy",
			"PolicyName":    "vcr_test_dl_policy",
			"PolicyVhost":   "/",
			"PolicyPattern": "^dl_test",
		}

		paramsUpdated = map[string]string{
			"ResourceName":         "dead_letter_policy",
			"PolicyName":           "vcr_test_dl_policy",
			"PolicyVhost":          "/",
			"PolicyPattern":        "^dl_test",
			"DeadLetterExchange":   "updated_dlx_exchange",
			"DeadLetterRoutingKey": "updated_dlx_routing_key",
		}
	)

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyResourceName, "name", params["PolicyName"]),
					resource.TestCheckResourceAttr(policyResourceName, "vhost", params["PolicyVhost"]),
					resource.TestCheckResourceAttr(policyResourceName, "pattern", params["PolicyPattern"]),
					resource.TestCheckResourceAttr(policyResourceName, "apply_to", "all"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyResourceName, "name", paramsUpdated["PolicyName"]),
					resource.TestCheckResourceAttr(policyResourceName, "vhost", paramsUpdated["PolicyVhost"]),
					resource.TestCheckResourceAttr(policyResourceName, "pattern", paramsUpdated["PolicyPattern"]),
					resource.TestCheckResourceAttr(policyResourceName, "definition.dead-letter-exchange", paramsUpdated["DeadLetterExchange"]),
					resource.TestCheckResourceAttr(policyResourceName, "definition.dead-letter-routing-key", paramsUpdated["DeadLetterRoutingKey"]),
					resource.TestCheckResourceAttr(policyResourceName, "apply_to", "all"),
				),
			},
		},
	})
}
