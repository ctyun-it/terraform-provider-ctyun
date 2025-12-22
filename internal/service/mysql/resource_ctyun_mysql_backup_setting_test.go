package mysql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunMysqlBackupSetting(t *testing.T) {

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_backup_setting." + rnd
	frequencyResourceFile := "resource_ctyun_mysql_backup_setting_frequency.tf"

	resourceFile := "resource_ctyun_mysql_backup_setting.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"
	mysqlInstanceID := dependence.mysqlID

	// 初始配置
	initialConfig := map[string]interface{}{
		"storage_day":                7,
		"frequency_backup":           false,
		"frequency_backup_unit_time": 3600,
		"allow_earliest_time":        "00:01",
		"trigger_days_of_week":       `[1, 3, 5]`, // 周一、三、五
	}

	// 更新配置
	updatedConfig := map[string]interface{}{
		"storage_day":                14,
		"frequency_backup":           true,
		"frequency_backup_unit_time": 7200, // 2小时
		"allow_earliest_time":        "02:00",
		"trigger_days_of_week":       `[2, 4, 6]`, // 周二、四、六
	}

	// 等待函数
	//wait10Seconds := func() {
	//	t.Logf("等待10秒...")
	//	time.Sleep(10 * time.Second)
	//}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建备份设置（初始配置）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, mysqlInstanceID, projectID,
					initialConfig["storage_day"], initialConfig["frequency_backup"], initialConfig["frequency_backup_unit_time"],
					initialConfig["allow_earliest_time"], initialConfig["trigger_days_of_week"],
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", mysqlInstanceID),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "storage_day", "7"),
					resource.TestCheckResourceAttr(resourceName, "frequency_backup", "false"),
					resource.TestCheckResourceAttr(resourceName, "allow_earliest_time", "00:01"),
					resource.TestCheckResourceAttr(resourceName, "trigger_days_of_week.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "trigger_days_of_week.0", "1"),
					resource.TestCheckResourceAttr(resourceName, "trigger_days_of_week.1", "3"),
					resource.TestCheckResourceAttr(resourceName, "trigger_days_of_week.2", "5"),
				),
			},
			// 2. 更新备份设置
			{
				Config: utils.LoadTestCase(frequencyResourceFile, rnd, mysqlInstanceID, projectID,
					updatedConfig["storage_day"], updatedConfig["frequency_backup"], updatedConfig["frequency_backup_unit_time"],
					updatedConfig["allow_earliest_time"], updatedConfig["trigger_days_of_week"],
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", mysqlInstanceID),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "storage_day", "14"),
					resource.TestCheckResourceAttr(resourceName, "frequency_backup", "true"),
					resource.TestCheckResourceAttr(resourceName, "frequency_backup_unit_time", "7200"),
					resource.TestCheckResourceAttr(resourceName, "allow_earliest_time", "02:00"),
					resource.TestCheckResourceAttr(resourceName, "trigger_days_of_week.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "trigger_days_of_week.0", "2"),
					resource.TestCheckResourceAttr(resourceName, "trigger_days_of_week.1", "4"),
					resource.TestCheckResourceAttr(resourceName, "trigger_days_of_week.2", "6"),
				),
			},
			// 4. 导入测试1
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					id := rs.Primary.Attributes["instance_id"]
					regionId := rs.Primary.Attributes["region_id"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,0,%s", id, regionId), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"}, // 不需要忽略任何字段
			},
			// 4. 导入测试2
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					id := rs.Primary.Attributes["instance_id"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s", id), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"}, // 不需要忽略任何字段
			},
			{
				Config: utils.LoadTestCase(frequencyResourceFile, rnd, mysqlInstanceID, projectID,
					updatedConfig["storage_day"], updatedConfig["frequency_backup"], updatedConfig["frequency_backup_unit_time"],
					updatedConfig["allow_earliest_time"], updatedConfig["trigger_days_of_week"],
				),
				Destroy: true,
			},
		},
	})
}
