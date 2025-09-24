package lavinmq

import (
	"testing"

	"github.com/cloudamqp/terraform-provider-lavinmq/lavinmq/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPolicy_Basic(t *testing.T) {
	var (
		fileNames          = []string{"policies/policy_basic"}
		policyResourceName = "lavinmq_policy.test_policy"

		params = map[string]string{
			"PolicyName":    "vcr_test_policy",
			"PolicyVhost":   "/",
			"PolicyPattern": "^vcr_test",
		}

		fileNamesUpdated_01 = []string{"policies/policy_with_definition"}
		paramsUpdated_01    = map[string]string{
			"PolicyName":      "vcr_test_policy",
			"PolicyVhost":     "/",
			"PolicyPattern":   "^vcr_test",
			"PolicyTTL":       "3600000",
			"PolicyMaxLength": "1000",
		}

		fileNamesUpdated_02 = []string{"policies/policy_with_priority"}
		paramsUpdated_02    = map[string]string{
			"PolicyName":     "vcr_test_policy",
			"PolicyVhost":    "/",
			"PolicyPattern":  "^vcr_test",
			"PolicyTTL":      "7200000",
			"PolicyPriority": "10",
		}

		fileNamesUpdated_03 = []string{"policies/policy_with_apply_to"}
		paramsUpdated_03    = map[string]string{
			"PolicyName":    "vcr_test_policy",
			"PolicyVhost":   "/",
			"PolicyPattern": "^vcr_test",
			"PolicyTTL":     "1800000",
			"PolicyApplyTo": "queues",
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
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesUpdated_01, paramsUpdated_01),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyResourceName, "name", paramsUpdated_01["PolicyName"]),
					resource.TestCheckResourceAttr(policyResourceName, "vhost", paramsUpdated_01["PolicyVhost"]),
					resource.TestCheckResourceAttr(policyResourceName, "pattern", paramsUpdated_01["PolicyPattern"]),
					resource.TestCheckResourceAttr(policyResourceName, "definition.message-ttl", paramsUpdated_01["PolicyTTL"]),
					resource.TestCheckResourceAttr(policyResourceName, "definition.max-length", paramsUpdated_01["PolicyMaxLength"]),
					resource.TestCheckResourceAttr(policyResourceName, "priority", "0"),
					resource.TestCheckResourceAttr(policyResourceName, "apply_to", "all"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesUpdated_02, paramsUpdated_02),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyResourceName, "name", paramsUpdated_02["PolicyName"]),
					resource.TestCheckResourceAttr(policyResourceName, "vhost", paramsUpdated_02["PolicyVhost"]),
					resource.TestCheckResourceAttr(policyResourceName, "pattern", paramsUpdated_02["PolicyPattern"]),
					resource.TestCheckResourceAttr(policyResourceName, "definition.message-ttl", paramsUpdated_02["PolicyTTL"]),
					resource.TestCheckResourceAttr(policyResourceName, "priority", paramsUpdated_02["PolicyPriority"]),
					resource.TestCheckResourceAttr(policyResourceName, "apply_to", "all"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesUpdated_03, paramsUpdated_03),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyResourceName, "name", paramsUpdated_03["PolicyName"]),
					resource.TestCheckResourceAttr(policyResourceName, "vhost", paramsUpdated_03["PolicyVhost"]),
					resource.TestCheckResourceAttr(policyResourceName, "pattern", paramsUpdated_03["PolicyPattern"]),
					resource.TestCheckResourceAttr(policyResourceName, "definition.message-ttl", paramsUpdated_03["PolicyTTL"]),
					resource.TestCheckResourceAttr(policyResourceName, "apply_to", paramsUpdated_03["PolicyApplyTo"]),
				),
			},
		},
	})
}

func TestAccPolicy_DeadLetter(t *testing.T) {
	var (
		fileNames          = []string{"policies/policy_dead_letter"}
		policyResourceName = "lavinmq_policy.dead_letter_policy"

		params = map[string]string{
			"PolicyName":           "vcr_test_dl_policy",
			"PolicyVhost":          "/",
			"PolicyPattern":        "^dl_test",
			"DeadLetterExchange":   "dlx_exchange",
			"DeadLetterRoutingKey": "dlx_routing_key",
		}

		paramsUpdated = map[string]string{
			"PolicyName":           "vcr_test_dl_policy",
			"PolicyVhost":          "/",
			"PolicyPattern":        "^dl_test",
			"DeadLetterExchange":   "updated_dlx_exchange",
			"DeadLetterRoutingKey": "updated_dlx_routing_key",
			"PolicyTTL":            "5000",
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
					resource.TestCheckResourceAttr(policyResourceName, "definition.dead-letter-exchange", params["DeadLetterExchange"]),
					resource.TestCheckResourceAttr(policyResourceName, "definition.dead-letter-routing-key", params["DeadLetterRoutingKey"]),
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
					resource.TestCheckResourceAttr(policyResourceName, "definition.message-ttl", paramsUpdated["PolicyTTL"]),
					resource.TestCheckResourceAttr(policyResourceName, "apply_to", "all"),
				),
			},
		},
	})
}

func TestAccPolicy_QueueLength(t *testing.T) {
	var (
		fileNames          = []string{"policies/policy_queue_length"}
		policyResourceName = "lavinmq_policy.queue_length_policy"

		params = map[string]string{
			"PolicyName":      "vcr_test_length_policy",
			"PolicyVhost":     "/",
			"PolicyPattern":   "^length_test",
			"PolicyMaxLength": "500",
			"PolicyApplyTo":   "queues",
		}

		paramsUpdated = map[string]string{
			"PolicyName":           "vcr_test_length_policy",
			"PolicyVhost":          "/",
			"PolicyPattern":        "^length_test",
			"PolicyMaxLength":      "1000",
			"PolicyMaxLengthBytes": "1048576",
			"PolicyApplyTo":        "queues",
			"PolicyPriority":       "5",
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
					resource.TestCheckResourceAttr(policyResourceName, "definition.max-length", params["PolicyMaxLength"]),
					resource.TestCheckResourceAttr(policyResourceName, "apply_to", params["PolicyApplyTo"]),
					resource.TestCheckResourceAttr(policyResourceName, "priority", "0"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyResourceName, "name", paramsUpdated["PolicyName"]),
					resource.TestCheckResourceAttr(policyResourceName, "vhost", paramsUpdated["PolicyVhost"]),
					resource.TestCheckResourceAttr(policyResourceName, "pattern", paramsUpdated["PolicyPattern"]),
					resource.TestCheckResourceAttr(policyResourceName, "definition.max-length", paramsUpdated["PolicyMaxLength"]),
					resource.TestCheckResourceAttr(policyResourceName, "definition.max-length-bytes", paramsUpdated["PolicyMaxLengthBytes"]),
					resource.TestCheckResourceAttr(policyResourceName, "apply_to", paramsUpdated["PolicyApplyTo"]),
					resource.TestCheckResourceAttr(policyResourceName, "priority", paramsUpdated["PolicyPriority"]),
				),
			},
		},
	})
}
