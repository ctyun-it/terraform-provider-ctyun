package pgsql_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"


	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
	prodId := 10003011
	storageType := "SSD"
	backupStorageType := `backup_storage_type="SATA"`
	StorageSpace := 100
	name := "pgsql-" + utils.GenerateRandomString()
	//password := "VqOcfgJ6Nf2houSe5C9sxgM4ycExVK+F0bBZwBGdiy8DCVXoSyck0lPxw9XMRgHur2lQYenOJ5K/FxZ30qlwbKG3NfgNoPq+AXDeSDdycGTqa1TzLdGnYwAeC/hEa8pyUKS9LdlW7nnM1nGUvGCXkGdzJP8lbHCwonzazEnF3RI="
	password := "Kyk123=2" + utils.GenerateRandomString()
	updatedPassword := "Kyk123==" + utils.GenerateRandomString()
	caseCensitive := true
	flavorName := dependence.flavorName
	updatedFlavorName := dependence.flavorName2
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	if dependence.securityGroupID == "" || dependence.securityGroupID2 == "" || dependence.securityGroupID3 == "" {
		t.Skip("three security_group_id values are required")
	}
	securityGroupID := fmt.Sprintf("%s,%s", dependence.securityGroupID, dependence.securityGroupID2)
	updatedSecurityGroupID := fmt.Sprintf("%s,%s", dependence.securityGroupID, dependence.securityGroupID3)
	azName := dependence.azName
	azInfo := fmt.Sprintf(`availability_zone_info=[{"availability_zone_name":"%s", "availability_zone_count":1, "node_type":"master"}]`, azName)

	updatedName := "pgsql-new" + utils.GenerateRandomString()
	//updatedSecurityGroupID := dependence.securityGroupID2
	updatedProdID := 10003012
	updatedStorageSpace := 120
	updatedAzInfo := fmt.Sprintf(`availability_zone_info=[{"availability_zone_name":"%s", "availability_zone_count":1, "node_type":"slave"}]`, azName)
	updatedBackupStorageSpace := fmt.Sprintf(`backup_storage_space="%d"`, updatedStorageSpace)
	vip := fmt.Sprintf("192.168.1.%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(101)+100)
	appointVIP := fmt.Sprintf(`vip="%s"`, vip)

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
			// Step 1：创建显式AZ的按需单节点实例，同时验证双安全组、备份盘和指定VIP。
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, flavorName, prodId, storageType, StorageSpace, name, password, caseCensitive,
					vpcID, subnetID, securityGroupID, azInfo, `backup_storage_space=100`, "", "", backupStorageType, appointVIP),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", 10003011)),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "vip", vip)),
			},
			// Step 2：修改root密码和安全组，验证两个可变属性复用同一实例。
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, flavorName, prodId, storageType, StorageSpace, name, updatedPassword, caseCensitive,
					vpcID, subnetID, updatedSecurityGroupID, azInfo, `backup_storage_space=100`, "", "", backupStorageType, appointVIP),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", updatedSecurityGroupID),
				),
			},
			// Step 3：修改实例名称并升级计算规格，不改变节点拓扑和磁盘容量。
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, prodId, storageType, StorageSpace, updatedName, updatedPassword, caseCensitive,
					vpcID, subnetID, updatedSecurityGroupID, azInfo, `backup_storage_space=100`, "", "", backupStorageType, appointVIP),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", 10003011)),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", updatedSecurityGroupID),
					resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", StorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
				),
			},
			// Step 4：仅扩容备节点磁盘，从100GB提升到120GB。
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, prodId, storageType, StorageSpace, updatedName, updatedPassword, caseCensitive,
					vpcID, subnetID, updatedSecurityGroupID, azInfo, updatedBackupStorageSpace, "", "", backupStorageType, appointVIP),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", 10003011)),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", updatedSecurityGroupID),
					resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", StorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_space", fmt.Sprintf("%d", updatedStorageSpace)),
				),
			},
			// Step 5：仅扩容主节点磁盘，从100GB提升到120GB。
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, prodId, storageType, updatedStorageSpace, updatedName, updatedPassword, caseCensitive,
					vpcID, subnetID, updatedSecurityGroupID, azInfo, updatedBackupStorageSpace, "", "", backupStorageType, appointVIP),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", 10003011)),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", updatedSecurityGroupID),
					resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", updatedStorageSpace)),
				),
			},
			// Step 6：单节点扩容为一主一备，并在扩容完成后停止实例。
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, updatedProdID, storageType, updatedStorageSpace, updatedName, updatedPassword, caseCensitive,
					vpcID, subnetID, updatedSecurityGroupID, updatedAzInfo, updatedBackupStorageSpace, `running_control="stop"`, "", backupStorageType, appointVIP),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", 10003012)),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", updatedSecurityGroupID),
				),
			},
			// Step 7：启动上一阶段已停止的实例，节点、规格和磁盘保持不变。
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, updatedProdID, storageType, updatedStorageSpace, updatedName, updatedPassword, caseCensitive,
					vpcID, subnetID, updatedSecurityGroupID, updatedAzInfo, updatedBackupStorageSpace, `running_control="start"`, "", backupStorageType, appointVIP),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", 10003012)),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", updatedSecurityGroupID),
				),
			},
			// Step 8：重启实例，验证restart运行控制，不再触发其他变配。
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, updatedProdID, storageType, updatedStorageSpace, updatedName, updatedPassword, caseCensitive,
					vpcID, subnetID, updatedSecurityGroupID, updatedAzInfo, updatedBackupStorageSpace, `running_control="restart"`, "", backupStorageType, appointVIP),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", 10003012)),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", updatedSecurityGroupID),
				),
			},
			// Step 9：使用instance_id查询datasource，验证更新后的实例可被准确检索。
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, updatedProdID, storageType, updatedStorageSpace, updatedName, updatedPassword, caseCensitive,
					vpcID, subnetID, updatedSecurityGroupID, updatedAzInfo, updatedBackupStorageSpace, "", "", backupStorageType, appointVIP) +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf("instance_id=%s.id", resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "instances.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "instances.0.name", updatedName),
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(30 * time.Second)
							return nil
						},
					),
				),
			},
			// Step 10：使用“实例ID,区域ID”格式导入并校验远端状态。
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
					return fmt.Sprintf("%s,%s", id, regionID), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"flavor_name", "case_sensitive", "master_order_id", "password"},
			},
			// Step 11：仅使用实例ID导入，验证默认区域导入路径。
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return s.RootModule().Resources[resourceName].Primary.ID, nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"flavor_name", "case_sensitive", "master_order_id", "password"},
			},
			// Step 12：销毁完成全部更新、运行控制、datasource和import验证的实例。
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, updatedFlavorName, updatedProdID, storageType, updatedStorageSpace, updatedName, updatedPassword, caseCensitive,
					vpcID, subnetID, updatedSecurityGroupID, updatedAzInfo, updatedBackupStorageSpace, "", "", backupStorageType, appointVIP),
				Destroy: true,
			},
		},
	},
	)
}

// 不传az Info 测试
func TestAccCtyunPgsqlInstanceNoAZInfo(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_postgresql_instance." + rnd

	resourceFile := "resource_ctyun_pgsql_instance.tf"

	cycleType := "on_demand"
	flavorName := dependence.flavorName
	prodId := 10003013
	storageType := "SSD"
	backupStorageType := `backup_storage_type = "SATA"`
	storageSpace := 100
	name := "pgsql-" + utils.GenerateRandomString()
	password := "Kyk123=2" + utils.GenerateRandomString()
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	backupStorageSpace := `backup_storage_space=100`

	updatedProdId := 10003014
	updatedStorageSpace := 150
	updatedBackupStorageSpace := `backup_storage_space = 200`
	updatedFlavorName := dependence.flavorName2
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
			// Step 1：不传AZ创建单节点实例，验证SATA备份盘及初始主备磁盘容量。
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, flavorName, prodId, storageType, storageSpace, name, password, caseCensitive,
					vpcID, subnetID, securityGroupID, "", backupStorageSpace, "", "", backupStorageType, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", prodId)),
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
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", updatedProdId)),
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
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", updatedProdId)),
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
	flavorName := dependence.flavorName
	prodId := 10003027
	storageType := "SSD"
	backupStorageType := `backup_storage_type="SAS"`
	storageSpace := 100
	name := "pgsql-" + utils.GenerateRandomString()
	password := "Kyk123=2" + utils.GenerateRandomString()
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	backupStorageSpace := `backup_storage_space=100`

	updatedStorageSpace := 150
	updatedBackupStorageSpace := `backup_storage_space=200`
	updatedFlavorName := dependence.flavorName2
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
			// Step 1：不传AZ直接创建一主两备实例，验证SAS备份盘和初始容量。
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, flavorName, prodId, storageType, storageSpace, name, password, caseCensitive,
					vpcID, subnetID, securityGroupID, "", backupStorageSpace, "", "", backupStorageType, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", prodId)),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", "SAS"),
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
					resource.TestCheckResourceAttr(resourceName, "prod_id", fmt.Sprintf("%d", prodId)),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", fmt.Sprintf("%d", updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", "SAS"),
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
