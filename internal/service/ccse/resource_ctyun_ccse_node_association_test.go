package ccse_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunCcseNodeAssociationEcs(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_ccse_node_association." + rnd
	resourceFile := "resource_ctyun_ccse_node_association.tf"

	clusterID := dependence.clusterID
	instanceType := "ecs"
	instanceID := dependence.ecsID
	mirrorID := dependence.ecsMirrorID
	visibilityPostHostScript := "YWJj"
	visibilityHostScript := "MTIz"
	password := "P@ss" + utils.GenerateRandomString()

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
				Config: utils.LoadTestCase(resourceFile, rnd,
					clusterID,
					instanceType,
					instanceID,
					mirrorID,
					visibilityPostHostScript,
					visibilityHostScript,
					password,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cluster_id", clusterID),
					resource.TestCheckResourceAttr(resourceName, "instance_type", instanceType),
					resource.TestCheckResourceAttr(resourceName, "instance_id", instanceID),
					resource.TestCheckResourceAttr(resourceName, "mirror_id", mirrorID),
					resource.TestCheckResourceAttr(resourceName, "visibility_post_host_script", visibilityPostHostScript),
					resource.TestCheckResourceAttr(resourceName, "visibility_host_script", visibilityHostScript),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "default_pool_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"az_name",
					"instance_id",
					"instance_type",
					"mirror_id",
					"visibility_host_script",
					"visibility_post_host_script",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					clusterID,
					instanceType,
					instanceID,
					mirrorID,
					visibilityPostHostScript,
					visibilityHostScript,
					password,
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunCcseNodeAssociationEbm(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_ccse_node_association." + rnd
	resourceFile := "resource_ctyun_ccse_node_association_ebm.tf"

	clusterID := dependence.clusterID
	instanceType := "ebm"
	instanceID := dependence.ebmID
	mirrorID := dependence.ebmMirrorID
	visibilityPostHostScript := "YWJj"
	visibilityHostScript := "MTIz"
	password := "P@ss" + utils.GenerateRandomString()

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
				Config: utils.LoadTestCase(resourceFile, rnd,
					clusterID,
					instanceType,
					instanceID,
					mirrorID,
					visibilityPostHostScript,
					visibilityHostScript,
					password,
					dependence.ebmAz,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cluster_id", clusterID),
					resource.TestCheckResourceAttr(resourceName, "instance_type", instanceType),
					resource.TestCheckResourceAttr(resourceName, "instance_id", instanceID),
					resource.TestCheckResourceAttr(resourceName, "mirror_id", mirrorID),
					resource.TestCheckResourceAttr(resourceName, "visibility_post_host_script", visibilityPostHostScript),
					resource.TestCheckResourceAttr(resourceName, "visibility_host_script", visibilityHostScript),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "default_pool_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"az_name",
					"instance_id",
					"instance_type",
					"mirror_id",
					"visibility_host_script",
					"visibility_post_host_script",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					clusterID,
					instanceType,
					instanceID,
					mirrorID,
					visibilityPostHostScript,
					visibilityHostScript,
					password,
					dependence.ebmAz,
				),
				Destroy: true,
			},
		},
	})
}
