package lavinmq

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccShovel_Import(t *testing.T) {
	t.Parallel()
	shovelResourceName := "lavinmq_shovel.test_shovel"

	// Set sanitized value for playback and use real value for recording
	testDestURI := "SHOVEL_DEST_URI"
	if os.Getenv("LAVINMQ_RECORD") != "" {
		testDestURI = os.Getenv("SHOVEL_DEST_URI")
	}

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
          resource "lavinmq_vhost" "test_vhost" {
            name	= "vcr-test-shovel-vhost-import"
          }
					
					resource "lavinmq_queue" "test_queue" {
						name       	= "vcr-test-shovel-queue-import"
						vhost      	= lavinmq_vhost.test_vhost.name
						durable    	= true
						auto_delete = false
					}
					
					resource "lavinmq_shovel" "test_shovel" {
						name               	= "vcr-test-shovel-import"
						vhost              	= lavinmq_vhost.test_vhost.name
						src_uri            	= "amqp://guest@/vcr-test-shovel-vhost-import"
						src_queue       		= lavinmq_queue.test_queue.name
						src_prefetch_count	= 1
						dest_uri    				= "%s"
						dest_queue  				= "vcr-test-shovel-destination-queue"
						reconnect_delay    	= 5
						ack_mode           	= "on-confirm"
					}
				`, testDestURI),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(shovelResourceName, "name", "vcr-test-shovel-import"),
					resource.TestCheckResourceAttr(shovelResourceName, "vhost", "vcr-test-shovel-vhost-import"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_uri", "amqp://guest@/vcr-test-shovel-vhost-import"),
					resource.TestCheckResourceAttr(shovelResourceName, "src_queue", "vcr-test-shovel-queue-import"),
					resource.TestCheckResourceAttr(shovelResourceName, "dest_queue", "vcr-test-shovel-destination-queue"),
					resource.TestCheckResourceAttr(shovelResourceName, "ack_mode", "on-confirm"),
				),
			},
			{
				ResourceName:                         shovelResourceName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateId:                        "vcr-test-shovel-vhost-import@vcr-test-shovel-import",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

func TestAccShovel_InvalidAttributeValue(t *testing.T) {
	t.Parallel()

	testNames := []string{"ack_mode", "src_delay_after"}
	configs := []string{
		`
			resource "lavinmq_shovel" "test_shovel" {
				name       = "vcr-test-shovel"
				vhost      = "vcr-test-shovel-vhost"
				src_uri    = "amqp://guest@/vcr-test-shovel-vhost"
				src_queue  = "vcr-test-shovel-src-queue"
				dest_uri   = "amqp://guest@/vcr-test-shovel-vhost"
				dest_queue = "vcr-test-shovel-dest-queue"
				ack_mode   = "invalid-mode"
			}
		`,
		`
			resource "lavinmq_shovel" "test_shovel" {
				name       			= "vcr-test-shovel"
				vhost      			= "vcr-test-shovel-vhost"
				src_uri    			= "amqp://guest@/vcr-test-shovel-vhost"
				src_queue  			= "vcr-test-shovel-src-queue"
				dest_uri   			= "amqp://guest@/vcr-test-shovel-vhost"
				dest_queue 			= "vcr-test-shovel-dest-queue"
				src_delay_after = "invalid-value"
			}
		`,
	}

	for index, config := range configs {
		t.Run(testNames[index], func(t *testing.T) {
			lavinMQResourceTest(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      config,
						ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
					},
				},
			})
		})
	}
}
