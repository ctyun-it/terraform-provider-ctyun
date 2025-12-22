package mysql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunMysqlReadOnlyInstance(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_readonly_instance." + rnd
	resourceFile := "resource_ctyun_mysql_read_node.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"
	instanceID := dependence.mysqlID
	cycleType := "on_demand"
	//cycleCount := 1
	// 测试数据
	instanceName := "test-ro-" + rnd
	flavorName := "c7.large.2" // 示例规格，根据实际情况调整
	storageType := "SATA"      // 存储类型
	storageSpace := 100        // 存储空间(GB)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建按需计费只读实例测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					instanceID, cycleType, flavorName,
					projectID, storageType, storageSpace, instanceName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cycle_type", "on_demand"),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			//// 2. 创建包月计费只读实例测试
			//{
			//	Config: testAccCtyunMysqlReadOnlyInstanceConfig(
			//		rnd, resourceFile,
			//		masterInstanceID, "month",
			//		1, true, flavorName,
			//		regionID, projectID, storageType,
			//		storageSpace, instanceName, availabilityZone,
			//	),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr(resourceName, "cycle_type", "month"),
			//		resource.TestCheckResourceAttr(resourceName, "cycle_count", "1"),
			//		resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
			//		resource.TestCheckResourceAttrSet(resourceName, "id"),
			//	),
			//},
			//// 3. 更新存储空间测试
			//{
			//	Config: testAccCtyunMysqlReadOnlyInstanceConfig(
			//		rnd, resourceFile,
			//		masterInstanceID, "month",
			//		1, true, flavorName,
			//		regionID, projectID, storageType,
			//		storageSpace+50, // 增加存储空间
			//		instanceName, availabilityZone,
			//	),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", storageSpace+50)),
			//	),
			//},
			// 4. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					instanceID, cycleType, flavorName,
					projectID, storageType, storageSpace, instanceName,
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunMysqlReadOnlyInstanceImportState(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_readonly_instance." + rnd
	resourceFile := "resource_ctyun_mysql_read_node.tf"

	// 从环境变量获取测试依赖资源
	projectID := "0"
	instanceID := dependence.mysqlID
	cycleType := "on_demand"
	//cycleCount := 1
	// 测试数据
	instanceName := "test-ros-" + rnd
	flavorName := "c7.large.2" // 示例规格，根据实际情况调整
	storageType := "SATA"      // 存储类型
	storageSpace := 100        // 存储空间(GB)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建按需计费只读实例测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					instanceID, cycleType, flavorName,
					projectID, storageType, storageSpace, instanceName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cycle_type", "on_demand"),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
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
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"cycle_type", "cycle_count", "auto_renew", "flavor_name",
					"storage_type", "storage_space", "availability_zone_name", "instance_id"}, // 不需要忽略任何字段
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
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"cycle_type", "cycle_count", "auto_renew", "flavor_name",
					"storage_type", "storage_space", "availability_zone_name", "instance_id"}, // 不需要忽略任何字段
			},
			// 4. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					instanceID, cycleType, flavorName,
					projectID, storageType, storageSpace, instanceName,
				),
				Destroy: true,
			},
		},
	})
}
