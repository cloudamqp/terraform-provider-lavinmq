package lavinmq

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-lavinmq/lavinmq/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccExchange_Basic(t *testing.T) {
	var (
		fileNames            = []string{"exchanges/exchange_basic"}
		exchangeResourceName = "lavinmq_exchange.vcr_test"

		params = map[string]string{
			"ExchangeName":  "vcr_test_exchange",
			"ExchangeVhost": "/",
			"ExchangeType":  "direct",
		}
	)

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(exchangeResourceName, "name", params["ExchangeName"]),
					resource.TestCheckResourceAttr(exchangeResourceName, "vhost", params["ExchangeVhost"]),
					resource.TestCheckResourceAttr(exchangeResourceName, "type", params["ExchangeType"]),
					resource.TestCheckResourceAttr(exchangeResourceName, "auto_delete", "false"),
					resource.TestCheckResourceAttr(exchangeResourceName, "durable", "false"),
				),
			},
			{
				ResourceName:                         exchangeResourceName,
				ImportStateIdFunc:                    testAccExchangeImportStateIdFunc,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"id"},
			},
		},
	})
}

func TestAccExchange_VhostScenarios(t *testing.T) {
	var (
		fileNamesCustomVhost = []string{"exchanges/exchange_custom_vhost"}
		exchangeResourceName = "lavinmq_exchange.vcr_test"

		paramsCustomVhost = map[string]string{
			"ExchangeName": "vcr_test_custom_vhost",
			"CustomVhost":  "test_vhost",
			"ExchangeType": "direct",
		}
	)

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesCustomVhost, paramsCustomVhost),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(exchangeResourceName, "name", paramsCustomVhost["ExchangeName"]),
					resource.TestCheckResourceAttr(exchangeResourceName, "vhost", paramsCustomVhost["CustomVhost"]),
					resource.TestCheckResourceAttr(exchangeResourceName, "type", paramsCustomVhost["ExchangeType"]),
					resource.TestCheckResourceAttr(exchangeResourceName, "auto_delete", "false"),
					resource.TestCheckResourceAttr(exchangeResourceName, "durable", "false"),
				),
			},
		},
	})
}

func testAccExchangeImportStateIdFunc(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "lavinmq_exchange" {
			continue
		}
		return fmt.Sprintf("%s,%s", rs.Primary.Attributes["vhost"], rs.Primary.Attributes["name"]), nil
	}
	return "", fmt.Errorf("resource not found")
}

func TestAccExchange_AllTypes(t *testing.T) {
	exchangeTypes := []struct {
		name     string
		filename string
		typeStr  string
	}{
		{"Direct", "exchanges/exchange_direct", "direct"},
		{"Fanout", "exchanges/exchange_fanout", "fanout"},
		{"Topic", "exchanges/exchange_topic", "topic"},
		{"Headers", "exchanges/exchange_headers", "headers"},
	}

	for _, exchangeType := range exchangeTypes {
		t.Run(exchangeType.name, func(t *testing.T) {
			var (
				fileNames            = []string{exchangeType.filename}
				exchangeResourceName = "lavinmq_exchange.vcr_test"

				params = map[string]string{
					"ExchangeName":  fmt.Sprintf("vcr_test_%s", exchangeType.typeStr),
					"ExchangeVhost": "/",
				}
			)

			lavinMQResourceTest(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: configuration.GetTemplatedConfig(t, fileNames, params),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(exchangeResourceName, "name", params["ExchangeName"]),
							resource.TestCheckResourceAttr(exchangeResourceName, "vhost", params["ExchangeVhost"]),
							resource.TestCheckResourceAttr(exchangeResourceName, "type", exchangeType.typeStr),
							resource.TestCheckResourceAttr(exchangeResourceName, "auto_delete", "false"),
							resource.TestCheckResourceAttr(exchangeResourceName, "durable", "false"),
						),
					},
				},
			})
		})
	}
}

func TestAccExchange_BooleanAttributes(t *testing.T) {
	var (
		fileNamesDurable     = []string{"exchanges/exchange_durable"}
		fileNamesAutoDelete  = []string{"exchanges/exchange_auto_delete"}
		exchangeResourceName = "lavinmq_exchange.vcr_test"

		paramsDurable = map[string]string{
			"ExchangeName":  "vcr_test_durable",
			"ExchangeVhost": "/",
			"ExchangeType":  "direct",
		}

		paramsAutoDelete = map[string]string{
			"ExchangeName":  "vcr_test_auto_delete",
			"ExchangeVhost": "/",
			"ExchangeType":  "direct",
		}
	)

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesDurable, paramsDurable),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(exchangeResourceName, "name", paramsDurable["ExchangeName"]),
					resource.TestCheckResourceAttr(exchangeResourceName, "vhost", paramsDurable["ExchangeVhost"]),
					resource.TestCheckResourceAttr(exchangeResourceName, "type", paramsDurable["ExchangeType"]),
					resource.TestCheckResourceAttr(exchangeResourceName, "durable", "true"),
					resource.TestCheckResourceAttr(exchangeResourceName, "auto_delete", "false"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesAutoDelete, paramsAutoDelete),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(exchangeResourceName, "name", paramsAutoDelete["ExchangeName"]),
					resource.TestCheckResourceAttr(exchangeResourceName, "vhost", paramsAutoDelete["ExchangeVhost"]),
					resource.TestCheckResourceAttr(exchangeResourceName, "type", paramsAutoDelete["ExchangeType"]),
					resource.TestCheckResourceAttr(exchangeResourceName, "auto_delete", "true"),
					resource.TestCheckResourceAttr(exchangeResourceName, "durable", "false"),
				),
			},
		},
	})
}

func TestAccExchange_Drift(t *testing.T) {
	var (
		fileNames = []string{"exchanges/exchange_drift"}

		params = map[string]string{
			"ExchangeName":  "vcr_test_drift",
			"ExchangeVhost": "/",
			"ExchangeType":  "direct",
		}
	)

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             configuration.GetTemplatedConfig(t, fileNames, params),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccExchange_WithArguments(t *testing.T) {
	var (
		fileNames            = []string{"exchanges/exchange_with_arguments"}
		exchangeResourceName = "lavinmq_exchange.vcr_test"

		params = map[string]string{
			"ExchangeName":         "vcr_test_exchange_args",
			"ExchangeVhost":        "/",
			"ExchangeType":         "direct",
			"ExchangeAutoDelete":   "false",
			"ExchangeDurable":      "true",
			"ExchangeArguments":    "true",
			"ArgAlternateExchange": "alternate_exchange",
		}
	)

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(exchangeResourceName, "name", params["ExchangeName"]),
					resource.TestCheckResourceAttr(exchangeResourceName, "vhost", params["ExchangeVhost"]),
					resource.TestCheckResourceAttr(exchangeResourceName, "type", params["ExchangeType"]),
					resource.TestCheckResourceAttr(exchangeResourceName, "auto_delete", "false"),
					resource.TestCheckResourceAttr(exchangeResourceName, "durable", "true"),
					resource.TestCheckResourceAttr(exchangeResourceName, "arguments.alternate-exchange", "alternate_exchange"),
				),
			},
			{
				ResourceName:                         exchangeResourceName,
				ImportStateIdFunc:                    testAccExchangeImportStateIdFunc,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"id"},
			},
		},
	})
}
