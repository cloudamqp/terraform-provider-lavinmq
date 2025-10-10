package lavinmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUser_Password(t *testing.T) {
	userResourceName := "lavinmq_user.user"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_user" "user" {
            name     = "vcr_test"
            password = "test1234"
            tags     = []
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", "vcr_test"),
					resource.TestCheckResourceAttr(userResourceName, "tags.#", "0"),
				),
			},
			{
				Config: `
          resource "lavinmq_user" "user" {
            name     = "vcr_test"
            password = "test1234"
            tags     = ["monitoring"]
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", "vcr_test"),
					resource.TestCheckResourceAttr(userResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(userResourceName, "tags.0", "monitoring"),
				),
			},
		},
	})
}

func TestAccUser_PasswordHash(t *testing.T) {
	userResourceName := "lavinmq_user.user"

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
          resource "lavinmq_user" "user" {
            name          = "vcr_test"
            password_hash = "jLmWeXp2gxeOfXhXOaHHWCf5mg1vYkMsIZeziH5ecf1HywxL"
            tags          = []
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", "vcr_test"),
					resource.TestCheckResourceAttr(userResourceName, "password_hash", "jLmWeXp2gxeOfXhXOaHHWCf5mg1vYkMsIZeziH5ecf1HywxL"),
					resource.TestCheckResourceAttr(userResourceName, "tags.#", "0"),
				),
			},
			{
				Config: `
          resource "lavinmq_user" "user" {
            name          = "vcr_test"
            password_hash = "jLmWeXp2gxeOfXhXOaHHWCf5mg1vYkMsIZeziH5ecf1HywxL"
            tags          = ["monitoring"]
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", "vcr_test"),
					resource.TestCheckResourceAttr(userResourceName, "password_hash", "jLmWeXp2gxeOfXhXOaHHWCf5mg1vYkMsIZeziH5ecf1HywxL"),
					resource.TestCheckResourceAttr(userResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(userResourceName, "tags.0", "monitoring"),
				),
			},
		},
	})
}
