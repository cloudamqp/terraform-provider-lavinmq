package lavinmq

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUser_Password(t *testing.T) {
	t.Parallel()
	userResourceName := "lavinmq_user.user"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "lavinmq_user" "user" {
						name     = "vcr-test-user-password"
						password = "test1234"
						tags     = ["monitoring"]
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", "vcr-test-user-password"),
					resource.TestCheckResourceAttr(userResourceName, "password_version", "1"),
					resource.TestCheckResourceAttr(userResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(userResourceName, "tags.0", "monitoring"),
				),
			},
			{
				Config: `
					resource "lavinmq_user" "user" {
						name     = "vcr-test-user-password"
						password = "test12345"
						password_version = 2
						tags     = ["monitoring"]
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", "vcr-test-user-password"),
					resource.TestCheckResourceAttr(userResourceName, "password_version", "2"),
					resource.TestCheckResourceAttr(userResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(userResourceName, "tags.0", "monitoring"),
				),
			},
		},
	})
}

func TestAccUser_PasswordHash(t *testing.T) {
	t.Parallel()
	userResourceName := "lavinmq_user.user"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "lavinmq_user" "user" {
						name          = "vcr-test-user-passwordhash"
						password_hash = {
							value     = "qV573OrTCGnMVbOysrKR2Xs16kkHiZbzhCDvf5mzV7NyH+M/"
							algorithm = "sha256"
						}
						tags          = ["monitoring"]
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", "vcr-test-user-passwordhash"),
					resource.TestCheckResourceAttr(userResourceName, "password_version", "1"),
					resource.TestCheckResourceAttr(userResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(userResourceName, "tags.0", "monitoring"),
				),
			},
			{
				Config: `
					resource "lavinmq_user" "user" {
						name          = "vcr-test-user-passwordhash"
						password_hash = {
							value     = "c6xQEdMpUle9NihE3SV8xcpXZtC6/z57IVlB22d/yEVw545L"
							algorithm = "sha256"
						}
						password_version = 2
						tags          = ["monitoring"]
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", "vcr-test-user-passwordhash"),
					resource.TestCheckResourceAttr(userResourceName, "password_version", "2"),
					resource.TestCheckResourceAttr(userResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(userResourceName, "tags.0", "monitoring"),
				),
			},
		},
	})
}

func TestAccUser_NoPasswordOrHash(t *testing.T) {
	t.Parallel()

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "lavinmq_user" "user" {
						name     = "vcr-test-user-no-password-or-hash"
						tags     = []
					}`,
				ExpectError: regexp.MustCompile("Either 'password' or 'password_hash' must be specified to create a user."),
			},
		},
	})
}

func TestAccUser_WithTags(t *testing.T) {
	t.Parallel()
	userResourceName := "lavinmq_user.user"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "lavinmq_user" "user" {
						name     = "vcr-test-user-with-tags"
						password = "test1234"
						tags     = ["monitoring"]
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", "vcr-test-user-with-tags"),
					resource.TestCheckResourceAttr(userResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(userResourceName, "tags.0", "monitoring"),
				),
			},
			{
				Config: `
					resource "lavinmq_user" "user" {
						name     = "vcr-test-user-with-tags"
						password = "test1234"
						tags     = ["monitoring", "management"]
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", "vcr-test-user-with-tags"),
					resource.TestCheckResourceAttr(userResourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(userResourceName, "tags.0", "monitoring"),
					resource.TestCheckResourceAttr(userResourceName, "tags.1", "management"),
				),
			},
		},
	})
}

func TestAccUser_WithoutTags(t *testing.T) {
	t.Parallel()
	userResourceName := "lavinmq_user.user"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "lavinmq_user" "user" {
						name     = "vcr-test-user-without-tags"
						password = "test1234"
						tags     = []
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", "vcr-test-user-without-tags"),
					resource.TestCheckResourceAttr(userResourceName, "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccUser_InvalidTag(t *testing.T) {
	t.Parallel()

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "lavinmq_user" "user" {
					name     = "vcr-test-user-invalid-tag"
					password = "test1234"
					tags     = ["invalid tag!"]
				}`,
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
		},
	})
}

func TestAccUser_Import(t *testing.T) {
	t.Parallel()
	userResourceName := "lavinmq_user.user"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "lavinmq_user" "user" {
						name     = "vcr-test-user-import"
						password = "test1234"
						tags     = ["monitoring"]
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", "vcr-test-user-import"),
					resource.TestCheckResourceAttr(userResourceName, "password_version", "1"),
					resource.TestCheckResourceAttr(userResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(userResourceName, "tags.0", "monitoring"),
				),
			},
			{
				ResourceName:                         userResourceName,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateId:                        "vcr-test-user-import",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}
