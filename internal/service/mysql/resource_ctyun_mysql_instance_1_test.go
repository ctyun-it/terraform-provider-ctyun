package mysql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunMysqlInstance1(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_instance." + rnd

	resourceFile := "resource_ctyun_mysql_instance.tf"
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mysql" + utils.GenerateRandomString()
	password := "Kyk123." + utils.GenerateRandomString()
	//period := 1
	//autoRenewStatus := 0

	storageType := "SATA"
	storageSpace := 100
	flavorName := "c7.xlarge.2"
	updatedDiskAvailabilityZoneInfo := fmt.Sprintf(`availability_zone_info = [{"availability_zone_name":"%s","availability_zone_count":2,"node_type":"slave"}]`, dependence.azName)
	// 单节点
	ProdId := "Single57"
	// 一主两备
	updatedDoubleProId := "Master2Slave57"
	cycleBillMode := "on_demand"
	NodeOneAvailabilityZoneInfo := fmt.Sprintf(`availability_zone_info = [{"availability_zone_name":"%s","availability_zone_count":1,"node_type":"master"}]`, dependence.azName)
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
			// 按需，单节点，升级为1主2备, 关机，开启，重启验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, vpcID, subnetID, securityGroupID, name, password, "", "", flavorName, ProdId, "", storageType, storageSpace, NodeOneAvailabilityZoneInfo, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Single57"),
				),
			},
			// 升级1主2备
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, vpcID, subnetID, securityGroupID, name, password, "", "", flavorName, updatedDoubleProId, "", storageType, storageSpace, updatedDiskAvailabilityZoneInfo, "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Master2Slave57"),
				),
			},
			// 关机验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, vpcID, subnetID, securityGroupID, name, password, "", "", flavorName, updatedDoubleProId, "", storageType, storageSpace, updatedDiskAvailabilityZoneInfo, `running_control="freeze"`, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_running_status", fmt.Sprintf("%d", 0)),
					resource.TestCheckResourceAttr(resourceName, "prod_order_status", fmt.Sprintf("%d", 6)),
				),
			},
			// 开机验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, vpcID, subnetID, securityGroupID, name, password, "", "", flavorName, updatedDoubleProId, "", storageType, storageSpace, updatedDiskAvailabilityZoneInfo, `running_control="unfreeze"`, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_running_status", fmt.Sprintf("%d", 0)),
					resource.TestCheckResourceAttr(resourceName, "prod_order_status", fmt.Sprintf("%d", 0)),
				),
			},
			// 重启验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleBillMode, vpcID, subnetID, securityGroupID, name, password, "", "", flavorName, updatedDoubleProId, "", storageType, storageSpace, updatedDiskAvailabilityZoneInfo, `running_control="restart"`, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "inst_id"),
					resource.TestCheckResourceAttr(resourceName, "prod_running_status", fmt.Sprintf("%d", 0)),
					resource.TestCheckResourceAttr(resourceName, "prod_order_status", fmt.Sprintf("%d", 0)),
				),
			},
			// 销毁
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, cycleBillMode, vpcID, subnetID, securityGroupID, name, password, "", "", flavorName, updatedDoubleProId, "", storageType, storageSpace, updatedDiskAvailabilityZoneInfo, "", ""),
				Destroy: true,
			},
		},
	})
}
