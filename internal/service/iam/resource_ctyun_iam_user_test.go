package iam_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

func TestAccCtyunIamUser(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_iam_user." + rnd
	datasourceName := "data.ctyun_iam_users." + dnd
	resourceFile := "resource_ctyun_iam_user.tf"
	datasourceFile := "datasource_ctyun_iam_users.tf"
	email := utils.GenerateRandomString() + "@example.com"
	phone := "13812345678"
	name := utils.GenerateRandomString()
	password := "P@ss2" + utils.GenerateRandomString()[:3]
	description := "init"
	groupID := dependence.groupID

	updatedEmail := utils.GenerateRandomString() + "@example.com"
	updatedPhone := "17912345678"
	updatedName := utils.GenerateRandomString()
	updatedPassword := "P@ss2" + utils.GenerateRandomString()[:3]
	updatedDescription := "updated"
	updatedGroupID := dependence.groupID2

	resource.Test(t, resource.TestCase{
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: utils.LoadTestCase(resourceFile, rnd, email, phone, name, password, description, groupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "phone", phone),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckTypeSetElemAttr(resourceName, "user_group_ids.*", groupID),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedEmail, updatedPhone, updatedName, updatedPassword, updatedDescription, updatedGroupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "email", updatedEmail),
					resource.TestCheckResourceAttr(resourceName, "phone", updatedPhone),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "password", updatedPassword),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckTypeSetElemAttr(resourceName, "user_group_ids.*", updatedGroupID),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedEmail, updatedPhone, updatedName, updatedPassword, updatedDescription, updatedGroupID) +
					utils.LoadTestCase(datasourceFile, dnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							ds := s.RootModule().Resources[datasourceName].Primary
							count, err := strconv.Atoi(ds.Attributes["users.#"])
							if err != nil || count == 0 {
								return fmt.Errorf("users 无效: %v", ds.Attributes)
							}

							for i := 0; i < count; i++ {
								if ds.Attributes[fmt.Sprintf("users.%d.email", i)] == updatedEmail {
									return nil
								}
							}
							return fmt.Errorf("未找到目标元素")
						},
					)),
			},

			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s",
						rs.Primary.Attributes["id"],
					), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedEmail, updatedPhone, updatedName, updatedPassword, updatedDescription, updatedGroupID) +
					utils.LoadTestCase(datasourceFile, dnd),
				Destroy: true,
			},
		},
	},
	)
}
