package ccse_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunCcseNamespace(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_ccse_namespace." + rnd
	resourceFile := "resource_ctyun_ccse_namespace.tf"

	datasourceName := "data.ctyun_ccse_namespaces." + dnd
	datasourceFile := "datasource_ctyun_ccse_namespaces.tf"

	valuesYaml := `apiVersion: v1
kind: Namespace
metadata:
  name: example
`
	valuesYaml2 := `---
apiVersion: "v1"
kind: "Namespace"
metadata:
  labels:
    kubernetes.io/metadata.name: "example"
  managedFields:
  - apiVersion: "v1"
    fieldsType: "FieldsV1"
    fieldsV1:
      f:metadata:
        f:labels:
          ".": {}
          f:kubernetes.io/metadata.name: {}
    manager: "fabric8-kubernetes-client"
  name: "example"
spec:
  finalizers:
  - "kubernetes"
status:
  phase: "Active"
`

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
					valuesYaml,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cluster_id", dependence.clusterID),
					resource.TestCheckResourceAttr(resourceName, "values_yaml", valuesYaml),
					resource.TestCheckResourceAttr(resourceName, "namespace", "example"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "actual_config"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					dependence.clusterID,
					valuesYaml2,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cluster_id", dependence.clusterID),
					resource.TestCheckResourceAttr(resourceName, "values_yaml", valuesYaml2),
					resource.TestCheckResourceAttr(resourceName, "namespace", "example"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "actual_config"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					dependence.clusterID,
					valuesYaml2,
				) + utils.LoadTestCase(datasourceFile, dnd, dependence.clusterID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "values_yaml"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					namespace := ds.Attributes["namespace"]
					cluster_id := ds.Attributes["cluster_id"]
					regionId := ds.Attributes["region_id"]
					return fmt.Sprintf("%s,%s,%s", namespace, cluster_id, regionId), nil
				},
				ImportStateVerifyIgnore: []string{
					"values_yaml",
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					namespace := ds.Attributes["namespace"]
					cluster_id := ds.Attributes["cluster_id"]
					return fmt.Sprintf("%s,%s", namespace, cluster_id), nil
				},
				ImportStateVerifyIgnore: []string{
					"values_yaml",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					dependence.clusterID,
					valuesYaml2,
				) + utils.LoadTestCase(datasourceFile, dnd, dependence.clusterID),
				Destroy: true,
			},
		},
	})
}
