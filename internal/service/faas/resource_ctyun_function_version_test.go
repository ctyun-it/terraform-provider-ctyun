package faas_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunFunctionVersion_Basic(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_function_version." + rnd
	//dataSourceName := "data.ctyun_function_versions." + rnd

	resourceFile := "resource_ctyun_function_version.tf"
	dataSourceFile := "datasource_ctyun_function_versions.tf"

	description := "test-version-" + rnd
	functionName := dependence.functionName

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
				// Step 1: 创建函数版本
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					description,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "function_name", functionName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "version_id"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				// Step 2: 导入测试
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					function_name := ds.Attributes["function_name"]
					version_id := ds.Attributes["version_id"]
					return fmt.Sprintf("%s,%s", function_name, version_id), nil
				},
			},
			{
				// Step 3: 数据源查询测试
				Config: utils.LoadTestCase(dataSourceFile,
					rnd,
					functionName,
				) + "\n" + utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					description,
				),
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
			{
				// Step 4: 销毁测试
				Config:  utils.LoadTestCase(resourceFile, rnd, functionName, description),
				Destroy: true,
			},
		},
	})
}
