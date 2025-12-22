package iam_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunIamPolicyAssociationUser(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_iam_policy_association_user." + rnd
	resourceFile := "resource_ctyun_iam_policy_association_user.tf"

	userID := dependence.userID
	policyID := dependence.policyID
	policyID2 := dependence.policyID2
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
				Config: utils.LoadTestCase(resourceFile, rnd, userID, policyID2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_id", userID),
					resource.TestCheckResourceAttr(resourceName, "policy_id", policyID2),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, userID, policyID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_id", userID),
					resource.TestCheckResourceAttr(resourceName, "policy_id", policyID),
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
					"region_id",
				},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, userID, policyID),
				Destroy: true,
			},
		},
	},
	)
}
