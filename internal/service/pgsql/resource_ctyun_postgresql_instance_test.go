package pgsql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
	"time"
)

func TestAccCtyunPgsqlInstance(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_postgresql_instance." + rnd
	datasourceName := "data.ctyun_postgresql_instances." + dnd

	resourceFile := "resource_ctyun_pgsql_instance.tf"
	datasourceFile := "datasource_ctyun_pgsql_instances.tf"

	cycleType := "on_demand"
	prodId := "Single1222"
	storageType := "SATA"
	backupStorageType := `backup_storage_type="SATA"`
	StorageSpace := 100
	name := "pgsql-" + utils.GenerateRandomString()
	//password := "VqOcfgJ6Nf2houSe5C9sxgM4ycExVK+F0bBZwBGdiy8DCVXoSyck0lPxw9XMRgHur2lQYenOJ5K/FxZ30qlwbKG3NfgNoPq+AXDeSDdycGTqa1TzLdGnYwAeC/hEa8pyUKS9LdlW7nnM1nGUvGCXkGdzJP8lbHCwonzazEnF3RI="
	password := "Kyk123="+utils.GenerateRandomString()
	caseCensitive := true
	flavorName := "s7.large.2"
	updatedFlavorName := "s7.large.4"
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	azName := dependence.azName
	azInfo := fmt.Sprintf(`availability_zone_info=[{"availability_zone_name":"%s", "availability_zone_count":1, "node_type":"master"}]`, azName)

	updatedName := "pgsql-new" + utils.GenerateRandomString()
	//updatedSecurityGroupID := dependence.securityGroupID2
	updatedProdID := "MasterSlave1222"
	updatedStorageSpace := 120
	updatedAzInfo := fmt.Sprintf(`availability_zone_info=[{"availability_zone_name":"%s", "availability_zone_count":1, "node_type":"slave"}]`, azName)
	updatedBackupStorageSpace := fmt.Sprintf(`backup_storage_space="%d"`, updatedStorageSpace)

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
					vpcID, subnetID, securityGroupID, azInfo, `backup_storage_space=100`, "", "", backupStorageType, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Single1222"),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName)),
			},
			// update验证--姓名, 安全组，规格扩容
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, prodId, storageType, StorageSpace, updatedName, password, caseCensitive,
					vpcID, subnetID, securityGroupID, azInfo, `backup_storage_space=100`, "", "", backupStorageType, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Single1222"),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", StorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
				),
			},
			// update验证--backup磁盘
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, prodId, storageType, StorageSpace, updatedName, password, caseCensitive,
					vpcID, subnetID, securityGroupID, azInfo, updatedBackupStorageSpace, "", "", backupStorageType, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Single1222"),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", StorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_space", fmt.Sprintf("%d", updatedStorageSpace)),
				),
			},
			// update验证--master磁盘
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, prodId, storageType, updatedStorageSpace, updatedName, password, caseCensitive,
					vpcID, subnetID, securityGroupID, azInfo, updatedBackupStorageSpace, "", "", backupStorageType, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Single1222"),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", updatedStorageSpace)),
				),
			},
			// update验证--主备，关机，开机，重启
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, updatedProdID, storageType, updatedStorageSpace, updatedName, password, caseCensitive,
					vpcID, subnetID, securityGroupID, updatedAzInfo, updatedBackupStorageSpace, `running_control="stop"`, "", backupStorageType, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "MasterSlave1222"),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, updatedProdID, storageType, updatedStorageSpace, updatedName, password, caseCensitive,
					vpcID, subnetID, securityGroupID, updatedAzInfo, updatedBackupStorageSpace, `running_control="start"`, "", backupStorageType, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "MasterSlave1222"),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, updatedProdID, storageType, updatedStorageSpace, updatedName, password, caseCensitive,
					vpcID, subnetID, securityGroupID, updatedAzInfo, "", `running_control="restart"`, "", backupStorageType, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "MasterSlave1222"),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			// datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, updatedProdID, storageType, updatedStorageSpace, updatedName, password, caseCensitive,
					vpcID, subnetID, securityGroupID, updatedAzInfo, "", ``, "", backupStorageType, "") +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf("prod_inst_id=%s.id", resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "pgsql_instances.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "pgsql_instances.0.name", updatedName),
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(30 * time.Second)
							return nil
						},
					),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, updatedProdID, storageType, updatedStorageSpace, updatedName, password, caseCensitive,
					vpcID, subnetID, securityGroupID, updatedAzInfo, "", ``, "", backupStorageType, ""),
				Destroy: true,
			},
		},
	},
	)
}

// 不传az Info 测试
func TestAccCtyunPgsqlInstanceNoAZInfo(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_postgresql_instance." + rnd

	resourceFile := "resource_ctyun_pgsql_instance.tf"

	cycleType := "on_demand"
	flavorName := "c7.large.2"
	prodId := "Single1417"
	storageType := "SAS"
	backupStorageType := `backup_storage_type = "SATA"`
	storageSpace := 100
	name := "pgsql-" + utils.GenerateRandomString()
	password := "Kyk123="+utils.GenerateRandomString()
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	backupStorageSpace := `backup_storage_space=100`

	updatedProdId := "MasterSlave1417"
	updatedStorageSpace := 150
	updatedBackupStorageSpace := `backup_storage_space = 200`
	updatedFlavorName := "c7.xlarge.2"
	caseCensitive := true

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
			{
				// 开通单结点
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, flavorName, prodId, storageType, storageSpace, name, password, caseCensitive,
					vpcID, subnetID, securityGroupID, "", backupStorageSpace, "", "", backupStorageType, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", prodId),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", "SATA"),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_space", "100")),
			},

			// 升级1主1备结点, 同时升级备份空间，主存储空间和spec
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, updatedProdId, storageType, storageSpace, name, password, caseCensitive,
					vpcID, subnetID, securityGroupID, "", backupStorageSpace, "", "", backupStorageType, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", updatedProdId),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName)),
			},
			// 升级1主1备结点, 同时升级备份空间，主存储空间和spec
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, updatedProdId, storageType, updatedStorageSpace, name, password, caseCensitive,
					vpcID, subnetID, securityGroupID, "", updatedBackupStorageSpace, "", "", backupStorageType, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", updatedProdId),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", "SATA"),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_space", "200"),
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(30 * time.Second)
							return nil
						},
					),
				),
			},
			// 销毁资源
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, updatedProdId, storageType, updatedStorageSpace, name, password, caseCensitive,
					vpcID, subnetID, securityGroupID, "", updatedBackupStorageSpace, "", "", backupStorageType, ""),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunPgsqlInstanceNoAZ2Info(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_postgresql_instance." + rnd

	resourceFile := "resource_ctyun_pgsql_instance.tf"

	cycleType := "on_demand"
	flavorName := "c7.large.2"
	prodId := "Master2Slave1512"
	storageType := "SSD"
	backupStorageType := `backup_storage_type="SSD"`
	storageSpace := 100
	name := "pgsql-" + utils.GenerateRandomString()
	password := "Kyk123="+utils.GenerateRandomString()
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	backupStorageSpace := `backup_storage_space=100`

	updatedStorageSpace := 150
	updatedBackupStorageSpace := `backup_storage_space=200`
	updatedFlavorName := "c7.large.4"
	caseCensitive := false
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
			{
				// 开通一主两备结点
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, flavorName, prodId, storageType, storageSpace, name, password, caseCensitive,
					vpcID, subnetID, securityGroupID, "", backupStorageSpace, "", "", backupStorageType, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", prodId),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", "SSD"),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_space", "100")),
			},
			// 升配主备磁盘，spec
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, prodId, storageType, updatedStorageSpace, name, password, caseCensitive,
					vpcID, subnetID, securityGroupID, "", updatedBackupStorageSpace, "", "", backupStorageType, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", prodId),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", "SSD"),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_space", "200"),
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(30 * time.Second)
							return nil
						},
					),
				),
			},

			// 销毁资源
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, prodId, storageType, updatedStorageSpace, name, password, caseCensitive,
					vpcID, subnetID, securityGroupID, "", updatedBackupStorageSpace, "", "", backupStorageType, ""),
				Destroy: true,
			},
		},
	})
}
