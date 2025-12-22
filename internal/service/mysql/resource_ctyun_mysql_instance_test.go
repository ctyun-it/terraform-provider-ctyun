package mysql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunMysqlInstanceUpdatePassword(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceFile := "resource_ctyun_mysql_instance_password.tf"
	resourceName := "ctyun_mysql_instance." + rnd
	cycleType := "on_demand"
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mysql-" + utils.GenerateRandomString()
	password := "Kyk111*" + utils.GenerateRandomString()
	prodID := "Single57"
	flavorName := "c7.xlarge.2"
	storageType := "SATA"
	storageSpace := 100
	updatePassword := "Kyk111**" + utils.GenerateRandomString()
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
			// 创建实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, name, password, flavorName, prodID, storageType, storageSpace),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 更新密码
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, name, updatePassword, flavorName, prodID, storageType, storageSpace),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// import state验证1
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					if id == "" || regionId == "" {
						return "", fmt.Errorf("id or region_id is required")
					}
					return fmt.Sprintf("%s,,%s", id, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"flavor_name", "password", "auto_renew",
					"backup_storage_type", "availability_zone_info", "running_control", "cycle_count", "master_order_id", "project_id"},
			},
			// import state验证2
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					return fmt.Sprintf("%s", id), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"flavor_name", "password", "auto_renew",
					"backup_storage_type", "availability_zone_info", "running_control", "master_order_id", "project_id"},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, name, password, flavorName, prodID, storageType, storageSpace),
				Destroy: true,
			},
		}})
}

func TestAccCtyunMysqlInstance(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_instance." + rnd
	datasourceName := "data.ctyun_mysql_instances." + dnd

	resourceFile := "resource_ctyun_mysql_instance.tf"
	datasourceFile := "datasource_ctyun_mysql_instances.tf"

	cycleType := "on_demand"
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mysql-" + utils.GenerateRandomString()
	password := "Kyk111*" + utils.GenerateRandomString()
	prodID := "Single57"
	flavorName := "c7.xlarge.2"

	storageType := "SATA"
	storageSpace := 100
	availabilityZoneInfo := fmt.Sprintf(`availability_zone_info = [{"availability_zone_name":"%s","availability_zone_count":1,"node_type":"master"}]`, dependence.azName)
	updatedDiskAvailabilityZoneInfo := fmt.Sprintf(`availability_zone_info = [{"availability_zone_name":"%s","availability_zone_count":1,"node_type":"slave"}]`, dependence.azName)
	updatedName := "tf-mysql-new-" + utils.GenerateRandomString()
	updatedWritePort := `write_port=13306`

	// 磁盘、规格升配
	updatedStorageSpace := 120
	updatedBackupStorageSpace := `backup_storage_space=150`
	updatedFlavorName := "c7.xlarge.4"
	// 单机到一主一备
	updatedProdID := "MasterSlave57"
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
			// 1. 按需验证，单节点创建，扩容至1主1备，修改端口，修改名称。
			// create 验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, name, password, "", "", flavorName, prodID, "", storageType, storageSpace, availabilityZoneInfo, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Single57"),
				),
			},
			// update, 实例名称、写端口更新验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, password, "", "", flavorName, prodID, updatedWritePort, storageType, storageSpace, availabilityZoneInfo, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Single57"),
					resource.TestCheckResourceAttr(resourceName, "write_port", "13306"),
				),
			},
			// 升配验证-升级磁盘空间
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, password, "", "", updatedFlavorName, prodID, updatedWritePort, storageType, updatedStorageSpace, availabilityZoneInfo, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "storage_space", "120"),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
				),
			},
			// 升配验证-升级备份磁盘空间
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, password, "", "", updatedFlavorName, prodID, updatedWritePort, storageType, updatedStorageSpace, availabilityZoneInfo, "", updatedBackupStorageSpace),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_space", "150"),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
				),
			},
			// 升配验证-单机规格扩容->1主1备
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, password, "", "", updatedFlavorName, updatedProdID, updatedWritePort, storageType, updatedStorageSpace, updatedDiskAvailabilityZoneInfo, "", updatedBackupStorageSpace),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "MasterSlave57"),
				),
			},
			// datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, password, "", "", updatedFlavorName, updatedProdID, updatedWritePort, storageType, updatedStorageSpace, updatedDiskAvailabilityZoneInfo, "", updatedBackupStorageSpace) +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf("name=%s.name", resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "instances.#", "1"),
				),
			},
			//销毁
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, password, "", "", updatedFlavorName, updatedProdID, updatedWritePort, storageType, updatedStorageSpace, updatedDiskAvailabilityZoneInfo, "", updatedBackupStorageSpace),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunMysqlInstanceMonth(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_instance." + rnd
	resourceFile := "resource_ctyun_mysql_instance.tf"
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mysql-" + utils.GenerateRandomString()
	password := "Kyk111*" + utils.GenerateRandomString()
	cycleCount := "cycle_count=1"
	autoRenewStatus := `auto_renew=false`

	flavorName := "c7.xlarge.2"

	storageType := "SATA"
	storageSpace := 100
	updatedDiskAvailabilityZoneInfo := fmt.Sprintf(`availability_zone_info = [{"availability_zone_name":"%s","availability_zone_count":1,"node_type":"slave"}]`, dependence.azName)

	// 单机到一主一备
	updatedProdID := "MasterSlave57"
	// 一主两备
	updatedDoubleProId := "Master2Slave57"
	cycleBillMode := "month"
	backupOneAvailabilityZoneInfo := fmt.Sprintf(`availability_zone_info=[{"availability_zone_name":"%s","availability_zone_count":1,"node_type":"master"},{"availability_zone_name":"%s","availability_zone_count":1,"node_type":"slave"}]`, dependence.azName, dependence.azName)

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
			// 2 包周期创建，创建1主1备，升级为1主2备,
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, vpcID, subnetID, securityGroupID, name, password, cycleCount, autoRenewStatus, flavorName, updatedProdID, "", storageType, storageSpace, backupOneAvailabilityZoneInfo, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "MasterSlave57"),
				),
			},
			// 升级1主2备
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, vpcID, subnetID, securityGroupID, name, password, cycleCount, autoRenewStatus, flavorName, updatedDoubleProId, "", storageType, storageSpace, updatedDiskAvailabilityZoneInfo, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Master2Slave57"),
				),
			},
			// 销毁
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, cycleBillMode, vpcID, subnetID, securityGroupID, name, password, cycleCount, autoRenewStatus, flavorName, updatedDoubleProId, "", storageType, storageSpace, updatedDiskAvailabilityZoneInfo, "", ""),
				Destroy: true,
			},
		},
	})
}
