package mysql_test

import (
	"fmt"
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

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
	updatedPassword := "Kyk111**" + utils.GenerateRandomString()
	prodID := "Single57"
	flavorName := dependence.flavorName

	storageType := "SATA"
	storageSpace := 100
	availabilityZoneInfo := fmt.Sprintf(`availability_zone_info = [{"availability_zone_name":"%s","availability_zone_count":1,"node_type":"master"}]`, dependence.azName)
	updatedDiskAvailabilityZoneInfo := fmt.Sprintf(`availability_zone_info = [{"availability_zone_name":"%s","availability_zone_count":1,"node_type":"slave"}]`, dependence.azName)
	updatedName := "tf-mysql-new-" + utils.GenerateRandomString()
	updatedWritePort := `write_port=13306`

	// 磁盘、规格升配
	updatedStorageSpace := 120
	updatedBackupStorageSpace := `backup_storage_space=150`
	updatedFlavorName := dependence.flavorName2
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
			// 更新root密码，并在后续步骤持续使用新密码，避免重复改密。
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, name, updatedPassword, "", "", flavorName, prodID, "", storageType, storageSpace, availabilityZoneInfo, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
				),
			},
			// update, 实例名称、写端口更新验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, updatedPassword, "", "", flavorName, prodID, updatedWritePort, storageType, storageSpace, availabilityZoneInfo, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Single57"),
					resource.TestCheckResourceAttr(resourceName, "write_port", "13306"),
				),
			},
			// 升配验证-升级磁盘空间
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, updatedPassword, "", "", updatedFlavorName, prodID, updatedWritePort, storageType, updatedStorageSpace, availabilityZoneInfo, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "storage_space", "120"),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
				),
			},
			// 升配验证-升级备份磁盘空间
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, updatedPassword, "", "", updatedFlavorName, prodID, updatedWritePort, storageType, updatedStorageSpace, availabilityZoneInfo, "", updatedBackupStorageSpace),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_space", "150"),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
				),
			},
			// 升配验证-单机规格扩容->1主1备
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, updatedPassword, "", "", updatedFlavorName, updatedProdID, updatedWritePort, storageType, updatedStorageSpace, updatedDiskAvailabilityZoneInfo, "", updatedBackupStorageSpace),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "MasterSlave57"),
				),
			},
			// datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, updatedPassword, "", "", updatedFlavorName, updatedProdID, updatedWritePort, storageType, updatedStorageSpace, updatedDiskAvailabilityZoneInfo, "", updatedBackupStorageSpace) +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf("name=%s.name", resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "instances.#", "1"),
				),
			},
			// import state验证：instance_id、region_id
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionID := ds.Attributes["region_id"]
					if id == "" || regionID == "" {
						return "", fmt.Errorf("id or region_id is required")
					}
					return fmt.Sprintf("%s,,%s", id, regionID), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"flavor_name", "password", "auto_renew",
					"backup_storage_type", "availability_zone_info", "running_control", "cycle_count", "master_order_id", "project_id", "prod_id"},
			},
			// import state验证：instance_id
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					return ds.ID, nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"flavor_name", "password", "auto_renew", "cycle_count",
					"backup_storage_type", "availability_zone_info", "running_control", "master_order_id", "project_id", "prod_id"},
			},
			//销毁
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, updatedPassword, "", "", updatedFlavorName, updatedProdID, updatedWritePort, storageType, updatedStorageSpace, updatedDiskAvailabilityZoneInfo, "", updatedBackupStorageSpace),
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

	flavorName := dependence.flavorName

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

// 使用数字格式的prod_id，覆盖显式AZ下的单机直升一主两备及运行控制。
func TestAccCtyunMysqlInstanceByInstanceType(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_instance." + rnd
	createFile := "resource_ctyun_mysql_instance_type_create.tf"
	updateFile := "resource_ctyun_mysql_instance_type_update.tf"

	cycleType := "on_demand"
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mysql-" + utils.GenerateRandomString()
	password := "Kyk111*" + utils.GenerateRandomString()
	prodID := "10001003"
	flavorName := dependence.flavorName

	storageType := "SSD-genric"
	storageSpace := 100
	availabilityZoneInfo := fmt.Sprintf(`[{"availability_zone_name":"%s","availability_zone_count":1,"node_type":"master"}]`, dependence.azName)
	updatedAvailabilityZoneInfo := fmt.Sprintf(`[{"availability_zone_name":"%s","availability_zone_count":2,"node_type":"slave"}]`, dependence.azName)
	updatedName := "tf-mysql-new-" + utils.GenerateRandomString()
	updatedWritePort := int32(13306)

	// 包周期参数（按需计费时 cycle_count 和 auto_renew 都必须为空）
	cycleCount := "" // on_demand 时 cycle_count 必须为空
	autoRenew := ""  // on_demand 时 auto_renew 必须为空

	// 磁盘、规格升配
	updatedStorageSpace := int32(110)
	updatedBackupStorageSpace := int32(110)
	updatedFlavorName := dependence.flavorName2
	// 单机到一主两备
	updatedProdID := "10001002"

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
			// 1. 按需验证，使用数字prod_id创建单节点。
			// create 验证
			{
				Config: utils.LoadTestCase(createFile, rnd, cycleType, vpcID, subnetID, securityGroupID, name, password, flavorName, prodID, storageType, storageSpace, availabilityZoneInfo),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "10001003"),
				),
			},
			// update, 实例名称、写端口更新验证
			{
				Config: utils.LoadTestCase(updateFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, password, cycleCount, autoRenew, flavorName, prodID, updatedWritePort, storageType, storageSpace, availabilityZoneInfo, storageSpace, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "10001003"),
					resource.TestCheckResourceAttr(resourceName, "write_port", "13306"),
				),
			},
			// 升配验证 - 升级磁盘空间
			{
				Config: utils.LoadTestCase(updateFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, password, cycleCount, autoRenew, updatedFlavorName, prodID, updatedWritePort, storageType, updatedStorageSpace, availabilityZoneInfo, storageSpace, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
				),
			},
			// 升配验证 - 升级备份磁盘空间
			{
				Config: utils.LoadTestCase(updateFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, password, cycleCount, autoRenew, updatedFlavorName, prodID, updatedWritePort, storageType, updatedStorageSpace, availabilityZoneInfo, updatedBackupStorageSpace, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
				),
			},
			// 升配验证 - 单机直升一主两备
			{
				Config: utils.LoadTestCase(updateFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, password, cycleCount, autoRenew, updatedFlavorName, updatedProdID, updatedWritePort, storageType, updatedStorageSpace, updatedAvailabilityZoneInfo, updatedBackupStorageSpace, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "10001002"),
				),
			},
			// 关机验证
			{
				Config: utils.LoadTestCase(updateFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, password, cycleCount, autoRenew, updatedFlavorName, updatedProdID, updatedWritePort, storageType, updatedStorageSpace, updatedAvailabilityZoneInfo, updatedBackupStorageSpace, `running_control="freeze"`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_running_status", "0"),
					resource.TestCheckResourceAttr(resourceName, "prod_order_status", "6"),
				),
			},
			// 开机验证
			{
				Config: utils.LoadTestCase(updateFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, password, cycleCount, autoRenew, updatedFlavorName, updatedProdID, updatedWritePort, storageType, updatedStorageSpace, updatedAvailabilityZoneInfo, updatedBackupStorageSpace, `running_control="unfreeze"`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_running_status", "0"),
					resource.TestCheckResourceAttr(resourceName, "prod_order_status", "0"),
				),
			},
			// 重启验证
			{
				Config: utils.LoadTestCase(updateFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, password, cycleCount, autoRenew, updatedFlavorName, updatedProdID, updatedWritePort, storageType, updatedStorageSpace, updatedAvailabilityZoneInfo, updatedBackupStorageSpace, `running_control="restart"`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_running_status", "0"),
					resource.TestCheckResourceAttr(resourceName, "prod_order_status", "0"),
				),
			},
			//销毁
			{
				Config:  utils.LoadTestCase(updateFile, rnd, cycleType, vpcID, subnetID, securityGroupID, updatedName, password, cycleCount, autoRenew, updatedFlavorName, updatedProdID, updatedWritePort, storageType, updatedStorageSpace, updatedAvailabilityZoneInfo, updatedBackupStorageSpace, ""),
				Destroy: true,
			},
		},
	})
}
