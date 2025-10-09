package lavinmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQueue_Import(t *testing.T) {
	queueResourceName := "lavinmq_queue.test_queue"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_queue" "test_queue" {
            name        = "vcr_test_queue_import"
            vhost       = "/"
            durable     = true
            auto_delete = false
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(queueResourceName, "name", "vcr_test_queue_import"),
					resource.TestCheckResourceAttr(queueResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(queueResourceName, "durable", "true"),
					resource.TestCheckResourceAttr(queueResourceName, "auto_delete", "false"),
					resource.TestCheckResourceAttr(queueResourceName, "pause", "false"),
					resource.TestCheckResourceAttrSet(queueResourceName, "state"),
				),
			},
			{
				ResourceName:      queueResourceName,
				ImportStateId:     "/,vcr_test_queue_import",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccQueue_PauseUnpause(t *testing.T) {
	queueResourceName := "lavinmq_queue.test_queue"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create queue in unpaused state
			{
				Config: `
          resource "lavinmq_queue" "test_queue" {
            name        = "vcr_test_queue_action"
            vhost       = "/"
            durable     = true
            auto_delete = false
            pause       = false
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(queueResourceName, "name", "vcr_test_queue_action"),
					resource.TestCheckResourceAttr(queueResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(queueResourceName, "durable", "true"),
					resource.TestCheckResourceAttr(queueResourceName, "pause", "false"),
					resource.TestCheckResourceAttr(queueResourceName, "state", "running"),
				),
			},
			// Step 2: Pause the queue
			{
				Config: `
          resource "lavinmq_queue" "test_queue" {
            name        = "vcr_test_queue_action"
            vhost       = "/"
            durable     = true
            auto_delete = false
            pause       = true
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(queueResourceName, "name", "vcr_test_queue_action"),
					resource.TestCheckResourceAttr(queueResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(queueResourceName, "durable", "true"),
					resource.TestCheckResourceAttr(queueResourceName, "pause", "true"),
					resource.TestCheckResourceAttr(queueResourceName, "state", "paused"),
				),
			},
			// Step 3: Unpause the queue
			{
				Config: `
          resource "lavinmq_queue" "test_queue" {
            name        = "vcr_test_queue_action"
            vhost       = "/"
            durable     = true
            auto_delete = false
            pause       = false
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(queueResourceName, "name", "vcr_test_queue_action"),
					resource.TestCheckResourceAttr(queueResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(queueResourceName, "durable", "true"),
					resource.TestCheckResourceAttr(queueResourceName, "pause", "false"),
					resource.TestCheckResourceAttr(queueResourceName, "state", "running"),
				),
			},
		},
	})
}
