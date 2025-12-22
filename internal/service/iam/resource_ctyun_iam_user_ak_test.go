package iam_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunIamUserAK(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_iam_user_ak." + rnd
	datasourceName := "data.ctyun_iam_user_aks." + dnd
	resourceFile := "resource_ctyun_iam_user_ak.tf"
	datasourceFile := "datasource_ctyun_iam_user_aks.tf"
	userID := dependence.userID

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
				Config: utils.LoadTestCase(resourceFile, rnd, userID, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_id", userID),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "ak"),
					resource.TestCheckResourceAttrSet(resourceName, "sk"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, userID, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_id", userID),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "ak"),
					resource.TestCheckResourceAttrSet(resourceName, "sk"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, userID, true) +
					utils.LoadTestCase(datasourceFile, dnd, userID, resourceName+".ak"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "user_id", userID),
					resource.TestCheckResourceAttr(datasourceName, "list.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "list.0.enabled", "true"),
					resource.TestCheckResourceAttrSet(datasourceName, "list.0.ak"),
					resource.TestCheckResourceAttrSet(datasourceName, "list.0.sk"),
					resource.TestCheckResourceAttrSet(datasourceName, "list.0.create_time"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s",
						rs.Primary.Attributes["ak"],
						rs.Primary.Attributes["user_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, userID, true) +
					utils.LoadTestCase(datasourceFile, dnd, userID, resourceName+".ak"),
				Destroy: true,
			},
		},
	},
	)
}
