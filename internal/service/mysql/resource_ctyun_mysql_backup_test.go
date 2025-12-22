package mysql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

// 创建备份集
func TestAccCtyunMysqlBackup(t *testing.T) {

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_backup." + rnd
	resourceFile := "resource_ctyun_mysql_backup.tf"

	dnd := utils.GenerateRandomString()
	backupsDatasourceName := "data.ctyun_mysql_backups." + dnd
	backupsDatasourceFile := "datasource_ctyun_mysql_backups.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"
	mysqlInstanceID := dependence.mysqlID
	// 备份描述信息
	description := "Test backup created by Terraform"

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
			// 1. 创建备份测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					mysqlInstanceID, projectID,
					description, "full",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "instance_id", mysqlInstanceID),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "task_type", "full"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
				),
			},
			// 2. 导入测试
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description", "task_type", "project_id"}, // 不需要忽略任何字段
			},
			// 验证 backups datasource
			{
				Config: utils.LoadTestCase(resourceFile, rnd, mysqlInstanceID, projectID, description, "full") +
					utils.LoadTestCase(backupsDatasourceFile, dnd, mysqlInstanceID, fmt.Sprintf("%s.name", resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(backupsDatasourceName, "backups.#"),
					resource.TestCheckResourceAttr(backupsDatasourceName, "backups.0.instance_id", mysqlInstanceID),
				),
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					mysqlInstanceID, projectID,
					description, "full",
				),
				Destroy: true,
			},
		},
	})
}

// import state
func TestAccCtyunMysqlBackupImportState(t *testing.T) {

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_backup." + rnd
	resourceFile := "resource_ctyun_mysql_backup.tf"

	projectID := "0"
	mysqlInstanceID := dependence.mysqlID
	// 备份描述信息
	description := "Test backup created by Terraform"

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
			// 1. 创建备份测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					mysqlInstanceID, projectID,
					description, "full",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "instance_id", mysqlInstanceID),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "task_type", "full"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
				),
			},
			// 3. 资源导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, _ := s.RootModule().Resources[resourceName]
					return fmt.Sprintf("%s,%s,%s,%s",
						rs.Primary.Attributes["name"],
						rs.Primary.Attributes["instance_id"],
						rs.Primary.Attributes["project_id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"task_type", "description"},
			},
			// 3. 资源导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, _ := s.RootModule().Resources[resourceName]
					return fmt.Sprintf("%s,%s",
						rs.Primary.Attributes["name"],
						rs.Primary.Attributes["instance_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"task_type", "description", "project_id"},
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					mysqlInstanceID, projectID,
					description, "full",
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunMysqlBackupCanceled(t *testing.T) {

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_backup." + rnd
	resourceFile := "resource_ctyun_mysql_backup.tf"

	cancelResourceName := "ctyun_mysql_backup_cancel." + rnd
	cancelResourceFile := "resource_ctyun_mysql_backup_cancel.tf"

	dnd := utils.GenerateRandomString()
	backupsDatasourceName := "data.ctyun_mysql_backups." + dnd
	backupsDatasourceFile := "datasource_ctyun_mysql_backups.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"
	mysqlInstanceID := dependence.mysqlID
	// 备份描述信息
	description := "Test backup created by Terraform"

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
			// 1. 创建备份测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					mysqlInstanceID, projectID,
					description, "full",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "instance_id", mysqlInstanceID),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "task_type", "full"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
				),
			},
			// datasource获取backup_record_id
			// 验证取消备份
			{
				Config: utils.LoadTestCase(resourceFile, rnd, mysqlInstanceID, projectID, description, "full") +
					utils.LoadTestCase(backupsDatasourceFile, dnd, mysqlInstanceID, fmt.Sprintf("%s.name", resourceName)) +
					utils.LoadTestCase(cancelResourceFile, rnd, mysqlInstanceID, projectID, fmt.Sprintf("%s.backups.0.records.0.backup_record_id", backupsDatasourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(cancelResourceName, "instance_id", mysqlInstanceID),
				),
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					mysqlInstanceID, projectID,
					description, "full",
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunMysqlBackupRecovery(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_backup_recovery." + rnd
	resourceFile := "resource_ctyun_mysql_backup_recovery.tf"

	timePointDatasourceName := "data.ctyun_mysql_recoverable_time_points." + dnd
	timePointDatasourceFile := "datasource_ctyun_mysql_recoverable_time_points.tf"

	backupResourceName := "ctyun_mysql_backup." + rnd
	backupResourceFile := "resource_ctyun_mysql_backup.tf"

	description := "Test backup created by Terraform"
	timeStamp := dependence.backupTimeStamp

	// 从环境变量获取测试依赖资源
	projectID := "0"
	srcInstanceID := dependence.mysqlID
	dstInstanceID := dependence.mysqlID
	//toTimepoint := os.Getenv("CTYUN_MYSQL_BACKUP_RECOVERY_TIME")

	// 等待函数
	//wait10Seconds := func() {
	//	t.Logf("等待10秒...")
	//	time.Sleep(10 * time.Second)
	//}

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
			// 1. 创建备份任务
			{
				Config: utils.LoadTestCase(
					backupResourceFile, rnd,
					srcInstanceID, projectID,
					description, "full",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(backupResourceName, "id"),
					resource.TestCheckResourceAttr(backupResourceName, "instance_id", srcInstanceID),
					resource.TestCheckResourceAttr(backupResourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(backupResourceName, "description", description),
					resource.TestCheckResourceAttr(backupResourceName, "task_type", "full"),
					resource.TestCheckResourceAttrSet(backupResourceName, "name"),
				),
			},
			// 2. 验证datasource
			{
				Config: utils.LoadTestCase(timePointDatasourceFile, dnd, srcInstanceID) +
					utils.LoadTestCase(
						backupResourceName, rnd,
						srcInstanceID, projectID,
						description, "full",
					),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(timePointDatasourceName, "backup_time_points.#"),
				),
			},
			// 3. 创建备份恢复任务
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, srcInstanceID, projectID, srcInstanceID, dstInstanceID,
					timeStamp,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", srcInstanceID),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "src_instance_id", srcInstanceID),
					resource.TestCheckResourceAttr(resourceName, "dst_instance_id", dstInstanceID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. 销毁
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, srcInstanceID, projectID, srcInstanceID, dstInstanceID,
					timeStamp,
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunMysqlBackupRecoveryByTaskID(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_backup_recovery." + rnd
	resourceFile := "resource_ctyun_mysql_backup_recovery_task_id.tf"

	backupsDatasourceName := "data.ctyun_mysql_backups." + dnd
	backupsDatasourceFile := "datasource_ctyun_mysql_backups_none_backup_name.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"
	srcInstanceID := dependence.mysqlID
	dstInstanceID := dependence.mysqlID

	taskID := dependence.taskID
	pageSize := 10
	pageNo := 1
	//toTimepoint := os.Getenv("CTYUN_MYSQL_BACKUP_RECOVERY_TIME")

	// 等待函数
	//wait10Seconds := func() {
	//	t.Logf("等待10秒...")
	//	time.Sleep(10 * time.Second)
	//}

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
			// 1. 验证datasource
			{
				Config: utils.LoadTestCase(backupsDatasourceFile, dnd, srcInstanceID, pageNo, pageSize),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(backupsDatasourceName, "backups.#"),
				),
			},
			// 1. 创建备份恢复任务
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, srcInstanceID, projectID, srcInstanceID, dstInstanceID, taskID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", srcInstanceID),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "src_instance_id", srcInstanceID),
					resource.TestCheckResourceAttr(resourceName, "dst_instance_id", dstInstanceID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. 销毁
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, srcInstanceID, projectID, srcInstanceID, dstInstanceID, taskID),
				Destroy: true,
			},
		},
	})
}
