package ccse_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunCcsePlugin(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_ccse_plugin." + rnd
	resourceFile := "resource_ctyun_ccse_plugin.tf"

	valuesYaml := fmt.Sprintf("values_yaml = %s", dependence.chartValuesYaml)
	valuesJson := fmt.Sprintf("values_json = %s", dependence.chartValuesJson)
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
					dependence.clusterID,
					dependence.chartName,
					dependence.chartVersion1,
					valuesYaml,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cluster_id", dependence.clusterID),
					resource.TestCheckResourceAttr(resourceName, "chart_name", dependence.chartName),
					resource.TestCheckResourceAttr(resourceName, "chart_version", dependence.chartVersion1),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					dependence.clusterID,
					dependence.chartName,
					dependence.chartVersion2,
					valuesYaml,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cluster_id", dependence.clusterID),
					resource.TestCheckResourceAttr(resourceName, "chart_name", dependence.chartName),
					resource.TestCheckResourceAttr(resourceName, "chart_version", dependence.chartVersion2),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					dependence.clusterID,
					dependence.chartName,
					dependence.chartVersion2,
					valuesJson,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cluster_id", dependence.clusterID),
					resource.TestCheckResourceAttr(resourceName, "chart_name", dependence.chartName),
					resource.TestCheckResourceAttr(resourceName, "chart_version", dependence.chartVersion2),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"values_json",
					"values_yaml",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					dependence.clusterID,
					dependence.chartName,
					dependence.chartVersion2,
					valuesJson,
				),
				Destroy: true,
			},
		},
	})
}
