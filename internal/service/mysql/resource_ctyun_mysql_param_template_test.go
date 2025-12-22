package mysql_test

import (
	"encoding/json"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunMysqlParamTemplate(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_param_template." + rnd
	createResourceFile := "resource_ctyun_mysql_param_template_create.tf"
	updateResourceFile := "resource_ctyun_mysql_param_template_update.tf"

	datasourceFile := "datasource_ctyun_mysql_param_templates.tf"
	datasourceName := "data.ctyun_mysql_param_templates." + dnd

	// 从环境变量获取测试依赖资源
	projectID := "0"
	engineVersion := "5.7"

	// 测试数据
	templateName := "test_template_" + rnd
	initialDescription := "Initial parameter template"
	updatedTemplateParameters := map[string]string{"auto_increment_increment": "2", "binlog_cache_size": "3000000"}
	updatedTemplateParametersStr, _ := json.Marshal(updatedTemplateParameters)
	//updatedDescription := "Updated parameter template"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建参数模板测试
			{
				Config: utils.LoadTestCase(
					createResourceFile, rnd, projectID,
					templateName, engineVersion, initialDescription,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "name", templateName),
					resource.TestCheckResourceAttr(resourceName, "engine", engineVersion),
					resource.TestCheckResourceAttr(resourceName, "description", initialDescription),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					//resource.TestCheckResourceAttrSet(resourceName, "template_parameters.#"),
				),
			},
			// 2. 更新描述信息测试
			{
				Config: utils.LoadTestCase(
					updateResourceFile, rnd, projectID,
					templateName, engineVersion, initialDescription, updatedTemplateParametersStr,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", initialDescription),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 4. 参数模板 datasource验证
			{
				Config: utils.LoadTestCase(datasourceFile, dnd, templateName, projectID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "param_templates.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "param_templates.0.name", templateName),
					resource.TestCheckResourceAttr(datasourceName, "param_templates.0.description", initialDescription),
				),
			},
			// 5. 清理资源
			{
				Config: utils.LoadTestCase(
					updateResourceFile, rnd, projectID,
					templateName, engineVersion, initialDescription, updatedTemplateParametersStr,
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunMysqlParamTemplateImportState(t *testing.T) {

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_param_template." + rnd
	createResourceFile := "resource_ctyun_mysql_param_template_create.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"
	engineVersion := "5.7"

	// 测试数据
	templateName := "test_template_" + rnd
	initialDescription := "Initial parameter template"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建参数模板测试
			{
				Config: utils.LoadTestCase(
					createResourceFile, rnd, projectID,
					templateName, engineVersion, initialDescription,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "name", templateName),
					resource.TestCheckResourceAttr(resourceName, "engine", engineVersion),
					resource.TestCheckResourceAttr(resourceName, "description", initialDescription),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					//resource.TestCheckResourceAttrSet(resourceName, "template_parameters.#"),
				),
			},
			// 2. 导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s,%s",
						rs.Primary.ID,
						rs.Primary.Attributes["project_id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "engine", "description", "template_parameters", "project_id"}, // 不需要忽略任何字段
			},
			// 3. 导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s",
						rs.Primary.ID,
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "engine", "description", "template_parameters", "project_id"}, // 不需要忽略任何字段
			},
			// 4. 清理资源
			{
				Config: utils.LoadTestCase(
					createResourceFile, rnd, projectID,
					templateName, engineVersion, initialDescription,
				),
				Destroy: true,
			},
		},
	})
}
