package pgsql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunPgsqlRunningControlInstance(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_postgresql_instance." + rnd

	resourceFile := "resource_ctyun_pgsql_instance.tf"

	cycleType := "on_demand"
	prodId := "Single1222"
	storageType := "SATA"
	backupStorageType := `backup_storage_type="SATA"`
	StorageSpace := 100
	name := "pgsql-" + utils.GenerateRandomString()
	//password := "VqOcfgJ6Nf2houSe5C9sxgM4ycExVK+F0bBZwBGdiy8DCVXoSyck0lPxw9XMRgHur2lQYenOJ5K/FxZ30qlwbKG3NfgNoPq+AXDeSDdycGTqa1TzLdGnYwAeC/hEa8pyUKS9LdlW7nnM1nGUvGCXkGdzJP8lbHCwonzazEnF3RI="
	password := "Kyk123=" + utils.GenerateRandomString()
	caseCensitive := true

	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	//azInfo := `availability_zone_info = [{"availability_zone_name":"cn-gs-qyi2-1a-public-ctcloud", "availability_zone_count":1, "node_type":"master"}]`

	flavorName := "s7.large.2"
	appointVip := `appoint_vip="192.168.4.111"`
	updatedFlavorName := "s7.large.4"
	updatedProdId := "MasterSlave1222"
	updatedStorageSpace := 120

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
			// 1. 按需验证，单节点创建，扩容至1主1备，修改名称，修改安全组， 规格扩容,磁盘扩容。
			// create 验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, flavorName, prodId, storageType, StorageSpace, name, password, caseCensitive,
					vpcID, subnetID, securityGroupID, "", `backup_storage_space=100`, "", "", backupStorageType, appointVip),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Single1222"),
					resource.TestCheckResourceAttr(resourceName, "appoint_vip", "192.168.4.111"),
				),
			},
			// 关机 + 主磁盘升配 + 备用磁盘升配 + sepc + prodid升配
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, updatedProdId, storageType, updatedStorageSpace, name, password, caseCensitive,
					vpcID, subnetID, securityGroupID, "", `backup_storage_space=120`, `running_control="stop"`, "", backupStorageType, appointVip),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "MasterSlave1222")),
			},
			// 开机
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, updatedProdId, storageType, updatedStorageSpace, name, password, caseCensitive,
					vpcID, subnetID, securityGroupID, "", `backup_storage_space=120`, `running_control="start"`, "", backupStorageType, appointVip),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "MasterSlave1222")),
			},
			// 重启
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, updatedProdId, storageType, updatedStorageSpace, name, password, caseCensitive,
					vpcID, subnetID, securityGroupID, "", `backup_storage_space=120`, `running_control="restart"`, "", backupStorageType, appointVip),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "MasterSlave1222")),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, updatedProdId, storageType, updatedStorageSpace, name, password, caseCensitive,
					vpcID, subnetID, securityGroupID, "", `backup_storage_space=120`, "", "", backupStorageType, appointVip),
				Destroy: true,
			},
		},
	},
	)
}
