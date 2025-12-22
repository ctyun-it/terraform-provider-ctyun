package redis_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunRedisMigrationTask(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_redis_migration_task." + rnd
	datasourceName := "data.ctyun_redis_migration_tasks." + dnd
	resourceFile := "resource_ctyun_redis_migration_task.tf"
	datasourceFile := "datasource_ctyun_redis_migration_tasks.tf"

	// 使用依赖中提供的Redis实例信息
	sourceInstanceId := dependence.instanceId
	sourceIp := dependence.address
	sourceAccount := dependence.userName
	sourcePassword := dependence.userPassword

	targetInstanceId := dependence.instance2Id
	targetIp := dependence.address2
	targetAccount := dependence.user2Name
	targetPassword := dependence.user2Password

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
			// 创建迁移任务
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					sourceInstanceId, sourceIp, sourceAccount, sourcePassword,
					targetInstanceId, targetIp, targetAccount, targetPassword, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "sync_mode", "1"),
					resource.TestCheckResourceAttr(resourceName, "conflict_mode", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			// 查询任务列表
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					sourceInstanceId, sourceIp, sourceAccount, sourcePassword,
					targetInstanceId, targetIp, targetAccount, targetPassword, "") +
					utils.LoadTestCase(datasourceFile, dnd, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "total"),
					resource.TestCheckResourceAttrSet(datasourceName, "size"),
					resource.TestCheckResourceAttrSet(datasourceName, "list.#"),
				),
			},
			// 查询在线迁移进度明细
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					sourceInstanceId, sourceIp, sourceAccount, sourcePassword,
					targetInstanceId, targetIp, targetAccount, targetPassword, "") +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf("id = ctyun_redis_migration_task.%s.id", rnd)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "source_progress_info_list.#"),
				),
			},
			// 结束运行中的任务
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					sourceInstanceId, sourceIp, sourceAccount, sourcePassword,
					targetInstanceId, targetIp, targetAccount, targetPassword, "operate_type = 2") +
					utils.LoadTestCase(datasourceFile, dnd, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "total"),
					resource.TestCheckResourceAttrSet(datasourceName, "size"),
					resource.TestCheckResourceAttrSet(datasourceName, "list.#"),
				),
			},

			// 导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs := s.RootModule().Resources[resourceName]
					if rs == nil {
						return "", fmt.Errorf("resource not found")
					}
					regionId := rs.Primary.Attributes["region_id"]
					taskId := rs.Primary.Attributes["id"]
					return fmt.Sprintf("%s,%s", taskId, regionId), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"operate_type", "source_db_info.password", "target_db_info.password"},
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs := s.RootModule().Resources[resourceName]
					if rs == nil {
						return "", fmt.Errorf("resource not found")
					}
					taskId := rs.Primary.Attributes["id"]
					return fmt.Sprintf("%s", taskId), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"operate_type", "source_db_info.password", "target_db_info.password"},
			},
			// 销毁测试
			{
				Config: utils.LoadTestCase(resourceFile, rnd,
					sourceInstanceId, sourceIp, sourceAccount, sourcePassword,
					targetInstanceId, targetIp, targetAccount, targetPassword, "") +
					utils.LoadTestCase(datasourceFile, dnd, ""),
				Destroy: true,
			},
		},
	})
}
