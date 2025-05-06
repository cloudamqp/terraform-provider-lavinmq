package lavinmq

import (
	"testing"

	"github.com/cloudamqp/terraform-provider-lavinmq/lavinmq/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVhost_Basic(t *testing.T) {
	var (
		fileNames         = []string{"vhosts/vhost_without_limits"}
		vhostResourceName = "lavinmq_vhost.vcr_test"

		params = map[string]string{
			"VhostName": "vcr_test",
		}

		fileNamesUpdated_01 = []string{"vhosts/vhost_with_limits"}
		paramsUpdated_01    = map[string]string{
			"VhostName":           "vcr_test",
			"VhostMaxConnections": "100",
			"VhostMaxQueues":      "30",
		}

		fileNamesUpdated_02 = []string{"vhosts/vhost_only_max_connections"}
		paramsUpdated_02    = map[string]string{
			"VhostName":           "vcr_test",
			"VhostMaxConnections": "100",
		}

		fileNamesUpdated_03 = []string{"vhosts/vhost_only_max_queues"}
		paramsUpdated_03    = map[string]string{
			"VhostName":      "vcr_test",
			"VhostMaxQueues": "30",
		}
	)

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(vhostResourceName, "name", params["VhostName"]),
					resource.TestCheckResourceAttr(vhostResourceName, "tracing", "false"),
					resource.TestCheckResourceAttr(vhostResourceName, "messages", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "messages_unacknowledged", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "messages_ready", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.ack", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.confirm", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.deliver", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.get", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.get_no_ack", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.publish", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.redeliver", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.return_unroutable", "0"),
					resource.TestCheckNoResourceAttr(vhostResourceName, "max_connections"),
					resource.TestCheckNoResourceAttr(vhostResourceName, "max_queues"),
				),
			},
			{
				ResourceName:                         vhostResourceName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateId:                        params["VhostName"],
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesUpdated_01, paramsUpdated_01),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(vhostResourceName, "name", paramsUpdated_01["VhostName"]),
					resource.TestCheckResourceAttr(vhostResourceName, "tracing", "false"),
					resource.TestCheckResourceAttr(vhostResourceName, "messages", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "messages_unacknowledged", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "messages_ready", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.ack", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.confirm", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.deliver", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.get", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.get_no_ack", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.publish", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.redeliver", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.return_unroutable", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "max_connections", paramsUpdated_01["VhostMaxConnections"]),
					resource.TestCheckResourceAttr(vhostResourceName, "max_queues", paramsUpdated_01["VhostMaxQueues"]),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesUpdated_02, paramsUpdated_02),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(vhostResourceName, "name", paramsUpdated_02["VhostName"]),
					resource.TestCheckResourceAttr(vhostResourceName, "messages", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "messages_unacknowledged", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "messages_ready", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.ack", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.confirm", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.deliver", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.get", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.get_no_ack", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.publish", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.redeliver", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.return_unroutable", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "max_connections", paramsUpdated_02["VhostMaxConnections"]),
					resource.TestCheckNoResourceAttr(vhostResourceName, "max_queues"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesUpdated_03, paramsUpdated_03),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(vhostResourceName, "name", paramsUpdated_03["VhostName"]),
					resource.TestCheckResourceAttr(vhostResourceName, "messages", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "messages_unacknowledged", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "messages_ready", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.ack", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.confirm", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.deliver", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.get", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.get_no_ack", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.publish", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.redeliver", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "message_stats.return_unroutable", "0"),
					resource.TestCheckResourceAttr(vhostResourceName, "max_queues", paramsUpdated_03["VhostMaxQueues"]),
					resource.TestCheckNoResourceAttr(vhostResourceName, "max_connections"),
				),
			},
		},
	})
}
