package pgsql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunPgsqlInstanceCycle(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_postgresql_instance." + rnd
	backupStorageType := `backup_storage_type="SSD"`
	resourceFile := "resource_ctyun_pgsql_instance.tf"
	cycleType := "month"
	prodId := "Single1222"
	storageType := "SATA"
	StorageSpace := 100
	name := "pgsql-" + utils.GenerateRandomString()
	//password := "VqOcfgJ6Nf2houSe5C9sxgM4ycExVK+F0bBZwBGdiy8DCVXoSyck0lPxw9XMRgHur2lQYenOJ5K/FxZ30qlwbKG3NfgNoPq+AXDeSDdycGTqa1TzLdGnYwAeC/hEa8pyUKS9LdlW7nnM1nGUvGCXkGdzJP8lbHCwonzazEnF3RI="
	password := "Kyk123=" + utils.GenerateRandomString()
	caseCensitive := true
	flavorName := "s7.large.2"
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	azName := dependence.azName
	azInfo := fmt.Sprintf(`availability_zone_info=[{"availability_zone_name":"%s", "availability_zone_count":1, "node_type":"master"}]`, azName)
	period := fmt.Sprint(`cycle_count=1`)

	updatedProdID := "Master2Slave1222"
	updatedAzInfo := fmt.Sprintf(`availability_zone_info=[{"availability_zone_name":"%s", "availability_zone_count":2, "node_type":"master"}]`, azName)

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
			// 按包周期创建单节点，并测试单节点->1主2备扩容
			// Create
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, flavorName, prodId, storageType, StorageSpace, name, password, caseCensitive,
					vpcID, subnetID, securityGroupID, azInfo, "", "", period, backupStorageType, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Single1222"),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", "SSD"),
				),
			},
			// 升配至1主2备
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, flavorName, updatedProdID, storageType, StorageSpace, name, password, caseCensitive,
					vpcID, subnetID, securityGroupID, updatedAzInfo, "", "", period, backupStorageType, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "prod_id", "Master2Slave1222"),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", "SSD"),
				),
			},
			// destroy
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, flavorName, updatedProdID, storageType, StorageSpace, name, password, caseCensitive,
					vpcID, subnetID, securityGroupID, updatedAzInfo, "", "", period, backupStorageType, ""),
				Destroy: true,
			},
		},
	})
}
