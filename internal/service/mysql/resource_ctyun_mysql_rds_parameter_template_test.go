package mysql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"strconv"
	"testing"
)

func TestAccCtyunMysqlRdsParameterTemplate(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_rds_parameter_template." + rnd
	resourceFile := "resource_ctyun_mysql_rds_parameter_template_bind.tf"
	updateResourceFile := "resource_ctyun_mysql_rds_parameter_template_update.tf"

	datasourceFile := "datasource_mysql_parameters.tf"
	datasourceName := "data.ctyun_mysql_parameters." + dnd

	// 从环境变量获取测试依赖资源
	projectID := "0"
	mysqlInstanceID := dependence.mysqlID
	templateID, err := strconv.Atoi(dependence.templateID)
	if err != nil {
		return
	}
	// 测试数据
	initialParameters := map[string]string{
		"auto_increment_increment": "65535",
		"automatic_sp_privileges":  "ON",
	}
	updatedParameters := map[string]string{
		"binlog_group_commit_sync_delay": "10",
		"back_log":                       "2000",
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// datasource验证
			{
				Config: utils.LoadTestCase(
					datasourceFile, dnd,
					projectID, mysqlInstanceID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "parameters.#")),
			},
			// 1. 应用参数模板测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					mysqlInstanceID, projectID,
					templateID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", mysqlInstanceID),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "template_id", fmt.Sprintf("%d", templateID)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. 直接修改参数测试
			{
				Config: utils.LoadTestCase(
					updateResourceFile, rnd,
					mysqlInstanceID, projectID,
					mapToTFConfigString(initialParameters),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "parameters.%", "2"),
				),
			},
			// 3. 更新参数测试
			{
				Config: utils.LoadTestCase(
					updateResourceFile, rnd,
					mysqlInstanceID, projectID,
					mapToTFConfigString(updatedParameters),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "parameters.%", "2"),
				),
			},
			// 4. 清理资源（恢复参数模板）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					mysqlInstanceID, projectID,
					templateID,
				),
				Destroy: true,
			},
		},
	})
}

func mapToTFConfigString(parameters map[string]string) string {
	if len(parameters) == 0 {
		return ""
	}

	paramsStr := ""
	for k, v := range parameters {
		if paramsStr != "" {
			paramsStr += ",\n"
		}
		paramsStr += fmt.Sprintf(`"%s" = "%s"`, k, v)
	}
	return fmt.Sprintf("{\n%s\n}", paramsStr)
}
