package lavinmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPermission_Basic(t *testing.T) {
	t.Parallel()
	var (
		permissionResourceName = "lavinmq_permission.test_permission"
	)

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_user" "test_user" {
            name     = "vcr_test_user"
            password = "test_password"
            tags     = ["management"]
          }

          resource "lavinmq_permission" "test_permission" {
            vhost     = "/"
            user      = lavinmq_user.test_user.name
            configure = ".*"
            read      = ".*"
            write     = ".*"
          }
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(permissionResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(permissionResourceName, "user", "vcr_test_user"),
					resource.TestCheckResourceAttr(permissionResourceName, "configure", ".*"),
					resource.TestCheckResourceAttr(permissionResourceName, "read", ".*"),
					resource.TestCheckResourceAttr(permissionResourceName, "write", ".*"),
				),
			},
		},
	})
}

func TestAccPermission_Update(t *testing.T) {
	t.Parallel()
	var (
		permissionResourceName = "lavinmq_permission.test_permission"
	)

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_user" "test_user" {
            name     = "vcr_test_user_update"
            password = "test_password"
            tags     = ["management"]
          }

          resource "lavinmq_permission" "test_permission" {
            vhost     = "/"
            user      = lavinmq_user.test_user.name
            configure = "^$"
            read      = ".*"
            write     = "^$"
          }
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(permissionResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(permissionResourceName, "user", "vcr_test_user_update"),
					resource.TestCheckResourceAttr(permissionResourceName, "configure", "^$"),
					resource.TestCheckResourceAttr(permissionResourceName, "read", ".*"),
					resource.TestCheckResourceAttr(permissionResourceName, "write", "^$"),
				),
			},
			{
				Config: `
          resource "lavinmq_user" "test_user" {
            name     = "vcr_test_user_update"
            password = "test_password"
            tags     = ["management"]
          }

          resource "lavinmq_permission" "test_permission" {
            vhost     = "/"
            user      = lavinmq_user.test_user.name
            configure = ".*"
            read      = ".*"
            write     = ".*"
          }
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(permissionResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(permissionResourceName, "user", "vcr_test_user_update"),
					resource.TestCheckResourceAttr(permissionResourceName, "configure", ".*"),
					resource.TestCheckResourceAttr(permissionResourceName, "read", ".*"),
					resource.TestCheckResourceAttr(permissionResourceName, "write", ".*"),
				),
			},
		},
	})
}

func TestAccPermission_Import(t *testing.T) {
	t.Parallel()
	var (
		permissionResourceName = "lavinmq_permission.test_permission"
	)

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_user" "test_user" {
            name     = "vcr_test_user_import"
            password = "test_password"
            tags     = ["management"]
          }

          resource "lavinmq_permission" "test_permission" {
            vhost     = "/"
            user      = lavinmq_user.test_user.name
            configure = ".*"
            read      = ".*"
            write     = ".*"
          }
        `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(permissionResourceName, "vhost", "/"),
					resource.TestCheckResourceAttr(permissionResourceName, "user", "vcr_test_user_import"),
					resource.TestCheckResourceAttr(permissionResourceName, "configure", ".*"),
					resource.TestCheckResourceAttr(permissionResourceName, "read", ".*"),
					resource.TestCheckResourceAttr(permissionResourceName, "write", ".*"),
				),
			},
			{
				ResourceName:                         permissionResourceName,
				ImportStateVerifyIdentifierAttribute: "user",
				ImportStateId:                        "/@vcr_test_user_import",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

func TestAccPermission_Drift(t *testing.T) {
	t.Parallel()
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_user" "test_user" {
            name     = "vcr_test_user_drift"
            password = "test_password"
            tags     = ["management"]
          }

          resource "lavinmq_permission" "test_permission" {
            vhost     = "/"
            user      = "vcr_test_user_drift"
            configure = ".*"
            read      = ".*"
            write     = ".*"
          }
        `,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
