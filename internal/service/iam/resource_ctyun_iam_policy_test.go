package iam_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunIamPolicy(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_iam_policy." + rnd
	datasourceName := "data.ctyun_iam_policies." + dnd
	resourceFile := "resource_ctyun_iam_policy.tf"
	datasourceFile := "datasource_ctyun_iam_policies.tf"

	name := "init"
	description := "init description"
	rAnge := "region"
	effect := "deny"
	action := dependence.authCode

	updatedName := "updated"
	updatedDescription := "updated description"
	updatedRange := "global"
	updatedEffect := "allow"
	updatedAction := dependence.authCode2

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
				Config: utils.LoadTestCase(resourceFile, rnd, name, description, rAnge, effect, action),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "range", rAnge),
					resource.TestCheckResourceAttr(resourceName, "content.statement.0.effect", effect),
					resource.TestCheckResourceAttr(resourceName, "content.statement.0.action.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "content.statement.0.action.*", action),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedDescription, updatedRange, updatedEffect, updatedAction),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttr(resourceName, "range", updatedRange),
					resource.TestCheckResourceAttr(resourceName, "content.statement.0.effect", updatedEffect),
					resource.TestCheckResourceAttr(resourceName, "content.statement.0.action.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "content.statement.0.action.*", updatedAction),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedDescription, updatedRange, updatedEffect, updatedAction) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "policies.#", "1"),
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
					return fmt.Sprintf("%s",
						rs.Primary.Attributes["id"],
					), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"action",
					"effect",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedDescription, updatedRange, updatedEffect, updatedAction) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Destroy: true,
			},
		},
	},
	)
}
