package lavinmq

import (
	"testing"

	"github.com/cloudamqp/terraform-provider-lavinmq/lavinmq/vcr-testing/configuration"
	"github.com/cloudamqp/terraform-provider-lavinmq/lavinmq/vcr-testing/converter"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUser_Password(t *testing.T) {
	t.Parallel()
	var (
		fileNames        = []string{"users/user_with_password"}
		userResourceName = "lavinmq_user.user"

		params = map[string]string{
			"UserName":     "vcr_test",
			"UserPassword": "test1234",
		}

		paramsUpdated = map[string]string{
			"UserName":     "vcr_test",
			"UserPassword": "test1234",
			"UserTags":     converter.CommaStringArray([]string{"monitoring"}),
		}
	)

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", params["UserName"]),
					resource.TestCheckResourceAttr(userResourceName, "tags.#", "0"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", paramsUpdated["UserName"]),
					resource.TestCheckResourceAttr(userResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(userResourceName, "tags.0", "monitoring"),
				),
			},
		},
	})
}

func TestAccUser_PasswordHash(t *testing.T) {
	t.Parallel()
	var (
		fileNames        = []string{"users/user_with_password_hash"}
		userResourceName = "lavinmq_user.user"

		params = map[string]string{
			"UserName":         "vcr_test",
			"UserPasswordHash": "jLmWeXp2gxeOfXhXOaHHWCf5mg1vYkMsIZeziH5ecf1HywxL",
		}

		paramsUpdated = map[string]string{
			"UserName":         "vcr_test",
			"UserPasswordHash": "jLmWeXp2gxeOfXhXOaHHWCf5mg1vYkMsIZeziH5ecf1HywxL",
			"UserTags":         converter.CommaStringArray([]string{"monitoring"}),
		}
	)

	lavinMQResourceTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", params["UserName"]),
					resource.TestCheckResourceAttr(userResourceName, "password_hash", params["UserPasswordHash"]),
					resource.TestCheckResourceAttr(userResourceName, "tags.#", "0"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", paramsUpdated["UserName"]),
					resource.TestCheckResourceAttr(userResourceName, "password_hash", paramsUpdated["UserPasswordHash"]),
					resource.TestCheckResourceAttr(userResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(userResourceName, "tags.0", "monitoring"),
				),
			},
		},
	})
}
