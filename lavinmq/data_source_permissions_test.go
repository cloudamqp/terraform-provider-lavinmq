package lavinmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourcePermissions_Basic(t *testing.T) {
	t.Parallel()
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_user" "test1" {
            name     = "terraform-lavinmq-test-permission-1"
            password = "test-password-1"
            tags     = ["management"]
          }

          resource "lavinmq_user" "test2" {
            name     = "terraform-lavinmq-test-permission-2"
            password = "test-password-2"
            tags     = ["management"]
          }

          resource "lavinmq_permission" "test1" {
            vhost     = "/"
            user      = lavinmq_user.test1.name
            configure = ".*"
            read      = ".*"
            write     = ".*"
          }

          resource "lavinmq_permission" "test2" {
            vhost     = "/"
            user      = lavinmq_user.test2.name
            configure = "^$"
            read      = ".*"
            write     = "^$"
          }

          data "lavinmq_permissions" "all" {
            depends_on = [lavinmq_permission.test1, lavinmq_permission.test2]
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_permissions.all", "permissions.*", map[string]string{
						"vhost":     "/",
						"user":      "terraform-lavinmq-test-permission-1",
						"configure": ".*",
						"read":      ".*",
						"write":     ".*",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_permissions.all", "permissions.*", map[string]string{
						"vhost":     "/",
						"user":      "terraform-lavinmq-test-permission-2",
						"configure": "^$",
						"read":      ".*",
						"write":     "^$",
					}),
				),
			},
		},
	})
}

func TestAccDataSourcePermissions_Empty(t *testing.T) {
	t.Parallel()
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          data "lavinmq_permissions" "all" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.lavinmq_permissions.all", "permissions.#"),
				),
			},
		},
	})
}

func TestAccDataSourcePermissions_FilterByVhost(t *testing.T) {
	t.Parallel()
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_vhost" "test_vhost" {
            name = "terraform-test-vhost"
          }

          resource "lavinmq_user" "test_user" {
            name     = "terraform-test-user"
            password = "test-password"
            tags     = ["management"]
          }

          resource "lavinmq_permission" "test" {
            vhost     = lavinmq_vhost.test_vhost.name
            user      = lavinmq_user.test_user.name
            configure = ".*"
            read      = ".*"
            write     = ".*"
          }

          data "lavinmq_permissions" "filtered" {
            vhost      = lavinmq_vhost.test_vhost.name
            depends_on = [lavinmq_permission.test]
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lavinmq_permissions.filtered", "vhost", "terraform-test-vhost"),
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_permissions.filtered", "permissions.*", map[string]string{
						"vhost":     "terraform-test-vhost",
						"user":      "terraform-test-user",
						"configure": ".*",
						"read":      ".*",
						"write":     ".*",
					}),
				),
			},
		},
	})
}

func TestAccDataSourcePermissions_FilterByUser(t *testing.T) {
	t.Parallel()
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_user" "test_user" {
            name     = "terraform-filter-user"
            password = "test-password"
            tags     = ["management"]
          }

          resource "lavinmq_permission" "test" {
            vhost     = "/"
            user      = lavinmq_user.test_user.name
            configure = "^test-.*"
            read      = ".*"
            write     = "^test-.*"
          }

          data "lavinmq_permissions" "filtered" {
            user       = lavinmq_user.test_user.name
            depends_on = [lavinmq_permission.test]
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lavinmq_permissions.filtered", "user", "terraform-filter-user"),
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_permissions.filtered", "permissions.*", map[string]string{
						"vhost":     "/",
						"user":      "terraform-filter-user",
						"configure": "^test-.*",
						"read":      ".*",
						"write":     "^test-.*",
					}),
				),
			},
		},
	})
}

func TestAccDataSourcePermissions_FilterByVhostAndUser(t *testing.T) {
	t.Parallel()
	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_vhost" "test_vhost" {
            name = "terraform-combo-vhost"
          }

          resource "lavinmq_user" "test_user" {
            name     = "terraform-combo-user"
            password = "test-password"
            tags     = ["management"]
          }

          resource "lavinmq_permission" "test" {
            vhost     = lavinmq_vhost.test_vhost.name
            user      = lavinmq_user.test_user.name
            configure = "^combo-.*"
            read      = ".*"
            write     = "^combo-.*"
          }

          data "lavinmq_permissions" "filtered" {
            vhost      = lavinmq_vhost.test_vhost.name
            user       = lavinmq_user.test_user.name
            depends_on = [lavinmq_permission.test]
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lavinmq_permissions.filtered", "vhost", "terraform-combo-vhost"),
					resource.TestCheckResourceAttr("data.lavinmq_permissions.filtered", "user", "terraform-combo-user"),
					resource.TestCheckResourceAttr("data.lavinmq_permissions.filtered", "permissions.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("data.lavinmq_permissions.filtered", "permissions.*", map[string]string{
						"vhost":     "terraform-combo-vhost",
						"user":      "terraform-combo-user",
						"configure": "^combo-.*",
						"read":      ".*",
						"write":     "^combo-.*",
					}),
				),
			},
		},
	})
}
