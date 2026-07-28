package mysql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunMysqlNoAzInfoInstance(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_instance." + rnd
	resourceFile := "resource_ctyun_mysql_instance.tf"

	cycleType := "on_demand"
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name1 := "tf-mysql-" + utils.GenerateRandomString()
	password := "Kyk111*" + utils.GenerateRandomString()
	updateProdID := "Master2Slave57"
	MsProdID := "MasterSlave57"

	storageType := "SATA"
	storageSpace := 100
	updatedStorageSpace := 120
	backupStorageSpace := `backup_storage_space=120`
	flavorName := dependence.flavorName

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
			// 直接开通一个1主1备的mysql, 并进行变配磁盘和1主2备
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, name1, password, "", "", flavorName, MsProdID, "", storageType, storageSpace, "", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name1),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "MasterSlave57"),
				),
			},
			// Step 2：不传AZ，将一主一备扩容为一主两备，同时扩容主备磁盘。
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, name1, password, "", "", flavorName, updateProdID, "", storageType, updatedStorageSpace, "", "", backupStorageSpace),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name1),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Master2Slave57"),
				),
			},
			// Step 3：销毁完成无AZ节点扩容和磁盘扩容验证的实例。
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, name1, password, "", "", flavorName, updateProdID, "", storageType, updatedStorageSpace, "", "", backupStorageSpace),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunMysqlNoAzInfoInstance1(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_instance." + rnd

	resourceFile := "resource_ctyun_mysql_instance.tf"

	cycleType := "on_demand"
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mysql-" + utils.GenerateRandomString()
	password := "Kyk111*" + utils.GenerateRandomString()
	prodID := "Single57"
	updateProdID := "Master2Slave57"

	storageType := "SATA"
	storageSpace := 100
	flavorName := dependence.flavorName
	updatedFlavorName := dependence.flavorName2

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
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, name, password, "", "", flavorName, prodID, "", storageType, storageSpace, "", "", ""),
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
			// Step 2：不传AZ，将单节点直接扩容为一主两备。
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, name, password, "", "", flavorName, updateProdID, "", storageType, storageSpace, "", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Master2Slave57"),
				),
			},
			// Step 3：保持一主两备拓扑不变，仅升级计算规格。
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, name, password, "", "", updatedFlavorName, updateProdID, "", storageType, storageSpace, "", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Master2Slave57"),
				),
			},
			// Step 4：销毁完成无AZ节点与规格升级验证的实例。
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, name, password, "", "", updatedFlavorName, updateProdID, "", storageType, storageSpace, "", "", ""),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunMysqlNoAzInfoInstance2(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_instance." + rnd

	resourceFile := "resource_ctyun_mysql_instance.tf"

	cycleType := "on_demand"
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID

	name2 := "tf-mysql-" + utils.GenerateRandomString()
	password := "Kyk111*" + utils.GenerateRandomString()
	updateProdID := "Master2Slave57"

	storageType := "SATA"
	storageSpace := 100
	flavorName := dependence.flavorName
	updatedFlavorName := dependence.flavorName2

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
			// 直接开通一个1主2备的mysql，变配规格
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, name2, password, "", "", flavorName, updateProdID, "", storageType, storageSpace, "", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name2),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Master2Slave57"),
				),
			},
			// 变配规格2c4g->2c8g
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, name2, password, "", "", updatedFlavorName, updateProdID, "", storageType, storageSpace, "", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name2),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Master2Slave57"),
				),
			},
			// Step 3：销毁直接创建并完成规格升级的一主两备实例。
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, subnetID, securityGroupID, name2, password, "", "", updatedFlavorName, updateProdID, "", storageType, storageSpace, "", "", ""),
				Destroy: true,
			},
		},
	})
}
