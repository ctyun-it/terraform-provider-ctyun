package pgsql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"strings"
	"testing"
)

func TestAccCtyunPostgresqlParamTemplate(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_postgresql_param_template." + rnd
	resourceFile := "resource_ctyun_postgresql_parameter_template.tf"
	resourceFile1 := "resource_ctyun_postgresql_parameter_template_update.tf"
	datasourceName := "data.ctyun_postgresql_param_templates." + dnd
	datasourceFile := "datasource_ctyun_postgresql_parameter_templates.tf"
	// 从环境变量获取测试依赖资源
	projectID := "0"
	sourceTemplateIDStr := dependence.paramTemplateID
	sourceTemplateID, err := strconv.Atoi(sourceTemplateIDStr)
	if err != nil {
		t.Fatal(err)
	}

	// 测试数据
	templateName := "tf-param-template-" + rnd
	initialDescription := "Initial_parameter_template_description"
	updatedDescription := "Updated_parameter_template_description"

	updatedParameters := map[string]string{
		"array_nulls":                     "off",
		"authentication_timeout":          "600",
		"autovacuum_analyze_scale_factor": "100",
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// datasource
			{
				Config: utils.LoadTestCase(datasourceFile, dnd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "parameter_templates.#"),
				),
			},
			// 1. 创建参数模板测试（基本配置）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					projectID,
					templateName,
					sourceTemplateID,
					initialDescription,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "name", templateName),
					resource.TestCheckResourceAttr(resourceName, "source_template_id", fmt.Sprintf("%d", sourceTemplateID)),
					resource.TestCheckResourceAttr(resourceName, "description", initialDescription),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckNoResourceAttr(resourceName, "template_parameters.%"),
				),
			},
			// 2. 更新参数模板描述测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					projectID,
					templateName,
					sourceTemplateID,
					updatedDescription,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
			},
			// 4. 更新参数测试
			{
				Config: utils.LoadTestCase(
					resourceFile1, rnd,
					projectID,
					templateName,
					sourceTemplateID,
					updatedDescription,
					MapToHCL(updatedParameters),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 5. 导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s,%s",
						rs.Primary.Attributes["id"],
						rs.Primary.Attributes["project_id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"template_parameters", "source_template_id", "description", "name"}, // 参数可能变化，忽略验证

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
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"template_parameters", "source_template_id", "description", "name", "project_id"}, // 参数可能变化，忽略验证

			},
			// 6. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					projectID,
					templateName,
					sourceTemplateID,
					updatedDescription,
				),
				Destroy: true,
			},
		},
	})
}

// MapToHCL 将 map 转换为 Terraform HCL 格式
func MapToHCL(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("{\n")
	for k, v := range m {
		builder.WriteString(fmt.Sprintf("  %s = \"%s\"\n", k, v))
	}
	builder.WriteString("}")
	return builder.String()
}
