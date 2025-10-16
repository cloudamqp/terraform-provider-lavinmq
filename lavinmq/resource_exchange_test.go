package lavinmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccExchange_Basic(t *testing.T) {
	t.Parallel()
	exchangeResourceName := "lavinmq_exchange.vcr_test"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_exchange" "vcr_test" {
            name        = "vcr_test_exchange"
            vhost       = "/"
            type        = "direct"
            auto_delete = false
            durable     = false
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(exchangeResourceName, "name", "vcr_test_exchange"),
					resource.TestCheckResourceAttr(exchangeResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(exchangeResourceName, "type", "direct"),
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
	t.Parallel()
	exchangeResourceName := "lavinmq_exchange.vcr_test"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_vhost" "test_vhost" {
            name = "test_vhost"
          }

          resource "lavinmq_exchange" "vcr_test" {
            name        = "vcr_test_custom_vhost"
            vhost       = lavinmq_vhost.test_vhost.name
            type        = "direct"
            auto_delete = false
            durable     = false
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(exchangeResourceName, "name", "vcr_test_custom_vhost"),
					resource.TestCheckResourceAttr(exchangeResourceName, "vhost", "test_vhost"),
					resource.TestCheckResourceAttr(exchangeResourceName, "type", "direct"),
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
	t.Parallel()
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
			exchangeResourceName := "lavinmq_exchange.vcr_test"
			exchangeName := fmt.Sprintf("vcr_test_%s", exchangeType.typeStr)

			lavinMQResourceTest(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(`
          resource "lavinmq_exchange" "vcr_test" {
            name        = "%s"
            vhost       = "/"
            type        = "%s"
            auto_delete = false
            durable     = false
          }`, exchangeName, exchangeType.typeStr),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(exchangeResourceName, "name", exchangeName),
							resource.TestCheckResourceAttr(exchangeResourceName, "vhost", "/"),
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
	t.Parallel()
	exchangeResourceName := "lavinmq_exchange.vcr_test"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_exchange" "vcr_test" {
            name        = "vcr_test_durable"
            vhost       = "/"
            type        = "direct"
            durable     = true
            auto_delete = false
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(exchangeResourceName, "name", "vcr_test_durable"),
					resource.TestCheckResourceAttr(exchangeResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(exchangeResourceName, "type", "direct"),
					resource.TestCheckResourceAttr(exchangeResourceName, "durable", "true"),
					resource.TestCheckResourceAttr(exchangeResourceName, "auto_delete", "false"),
				),
			},
			{
				Config: `
          resource "lavinmq_exchange" "vcr_test" {
            name        = "vcr_test_auto_delete"
            vhost       = "/"
            type        = "direct"
            auto_delete = true
            durable     = false
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(exchangeResourceName, "name", "vcr_test_auto_delete"),
					resource.TestCheckResourceAttr(exchangeResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(exchangeResourceName, "type", "direct"),
					resource.TestCheckResourceAttr(exchangeResourceName, "auto_delete", "true"),
					resource.TestCheckResourceAttr(exchangeResourceName, "durable", "false"),
				),
			},
		},
	})
}

func TestAccExchange_Drift(t *testing.T) {
	t.Parallel()
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_exchange" "vcr_test" {
            name        = "vcr_test_drift"
            vhost       = "/"
            type        = "direct"
            auto_delete = false
            durable     = false
          }`,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccExchange_WithArguments(t *testing.T) {
	t.Parallel()
	exchangeResourceName := "lavinmq_exchange.test_exchange"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_exchange" "test_exchange" {
            name        = "vcr_test_exchange_with_arguments"
            vhost       = "/"
            type        = "direct"
            durable     = true
            auto_delete = false
            arguments = {
              alternate-exchange = "alternate_exchange"
            }
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(exchangeResourceName, "name", "vcr_test_exchange_with_arguments"),
					resource.TestCheckResourceAttr(exchangeResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(exchangeResourceName, "type", "direct"),
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
