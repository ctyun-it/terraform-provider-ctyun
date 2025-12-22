package iam_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunEnterpriseProject(t *testing.T) {
	rnd := utils.GenerateRandomString()
	and := utils.GenerateRandomString()

	resourceName := "ctyun_enterprise_project." + rnd
	associationName := "ctyun_enterprise_project_association_user_group." + and
	resourceFile := "resource_ctyun_enterprise_project.tf"
	associationFile := "resource_ctyun_enterprise_project_association_user_group.tf"

	groupID := dependence.groupID
	policyID := dependence.policyID
	policyID2 := dependence.policyID2
	name := "init"
	description := "init description"

	updatedName := "updated"
	updatedDescription := "updated description"

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
				Config: utils.LoadTestCase(resourceFile, rnd, name, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedDescription) +
					utils.LoadTestCase(associationFile, and, resourceName+".id", groupID, policyID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(associationName, "user_group_id", groupID),
					resource.TestCheckTypeSetElemAttr(associationName, "policy_ids.*", policyID),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedDescription) +
					utils.LoadTestCase(associationFile, and, resourceName+".id", groupID, policyID2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(associationName, "user_group_id", groupID),
					resource.TestCheckTypeSetElemAttr(associationName, "policy_ids.*", policyID2),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, updatedDescription),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s",
						rs.Primary.Attributes["id"],
					), nil
				},
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, updatedName, updatedDescription),
				Destroy: true,
			},
		},
	},
	)
}
