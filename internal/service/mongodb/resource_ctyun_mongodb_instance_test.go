package mongodb_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
	"time"
)

// 单机、按需、有az、备份盘
func TestAccCtyunMongodbInstanceSingleOnDemand(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mongodb_instance." + rnd
	//datasourceName := "data.ctyun_mongodb_instances." + dnd

	resourceFile := "resource_ctyun_mongodb_instance_single_on_demand.tf"
	//datasourceFile := "datasource_ctyun_mongodb_instances.tf"
	datasourceName := "data.ctyun_mongodb_instances." + dnd
	datasourceFile := "datasource_ctyun_mongodb_instances.tf"
	// 创建参数
	cycleType := "on_demand"
	vpcID := dependence.vpcID
	flavorName := "s7.large.2"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mongodb-single-" + utils.GenerateRandomString()
	password := "Kyk123=" + utils.GenerateRandomString()
	prodId := "Single34"
	readPort := 12345
	storageType := "SAS"
	storageSpace := 120
	backupStorageType := "SATA"
	azName := dependence.azName
	azInfo := fmt.Sprintf(`[{"availability_zone_name":"%s","availability_zone_count":1,"node_type":"master"}, {"availability_zone_name":"%s","availability_zone_count":1,"node_type":"backup"}]`, azName, azName)

	//更新参数
	updatedName := "tf-mongodb-single-new-" + utils.GenerateRandomString()
	updatedFlavorName := "s7.large.4"
	updatedReadPort := 12348
	//updatedStorageType := ""
	updatedStorageSpace := 130
	//backupStorageType := "SATA"
	updatedAzInfo := fmt.Sprintf(`[{"availability_zone_name":"%s","availability_zone_count":1,"node_type":"s"}]`, azName)

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
			// 创建一个单节点的mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort, storageType, storageSpace,
					backupStorageType, azInfo),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			// 更新mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, prodId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedAzInfo),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(updatedReadPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			// datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, prodId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedAzInfo) +
					utils.LoadTestCase(datasourceFile, dnd, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "mongodb_instances.#"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, prodId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedAzInfo),
				Destroy: true,
			},
		},
	})
}

// 创建包周期，且无传AZ信息, 备份空间为os
func TestAccCtyunMongodbInstanceSingleCycleNoAz(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	//dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mongodb_instance." + rnd
	//datasourceName := "data.ctyun_mongodb_instances." + dnd

	resourceFile := "resource_ctyun_mongodb_instance_single_cycle_no_az_os.tf"
	//datasourceFile := "datasource_ctyun_mongodb_instances.tf"
	// 创建参数
	cycleType := "month"
	cycleCount := 1
	vpcID := dependence.vpcID
	flavorName := "s7.large.2"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mongodb-single-" + utils.GenerateRandomString()
	password := "Kyk123=" + utils.GenerateRandomString()
	prodId := "Single34"
	readPort := 12345
	storageType := "SATA"
	storageSpace := 100
	backupStorageType := "OS"
	//backupStorageSpace := 100
	//azInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"master"}, {"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"backup"}]`

	//更新参数
	updatedName := "tf-mongodb-single-new-" + utils.GenerateRandomString()
	updatedFlavorName := "s7.large.4"
	updatedReadPort := 12348
	//updatedStorageType := ""
	updatedStorageSpace := 110
	//backupStorageType := "SATA"
	//updatedBackupStorageSpace := 160

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
			// 创建一个单节点的mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, cycleCount, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort, storageType,
					storageSpace, backupStorageType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			// 更新mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, cycleCount, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, prodId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(updatedReadPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, cycleCount, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, prodId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType),
				Destroy: true,
			},
		},
	})
}

// 验证副本集，传azList，OS存储
func TestAccCtyunMongodbInstanceReplicaOs(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	//dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mongodb_instance." + rnd
	//datasourceName := "data.ctyun_mongodb_instances." + dnd

	resourceFile := "resource_ctyun_mongodb_instance_replica_on_demand_os.tf"
	//datasourceFile := "datasource_ctyun_mongodb_instances.tf"
	// 创建参数
	cycleType := "on_demand"
	vpcID := dependence.vpcID
	flavorName := "s7.large.2"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mongodb-single-" + utils.GenerateRandomString()
	password := "Kyk123=" + utils.GenerateRandomString()
	prodId := "Replica3R34"
	readPort := 12345
	storageType := "SAS"
	storageSpace := 100
	backupStorageType := "OS"
	azName := dependence.azName
	azInfo := fmt.Sprintf(`[{"availability_zone_name":"%s","availability_zone_count":3,"node_type":"master"}]`, azName)

	//更新参数
	updatedName := "tf-mongodb-single-new-" + utils.GenerateRandomString()
	updatedFlavorName := "s7.large.4"
	updatedReadPort := 12348
	updatedProdId := "Replica5R34"

	//updatedStorageType := ""
	updatedStorageSpace := 110
	updatedAzInfo := fmt.Sprintf(`[{"availability_zone_name":"%s","availability_zone_count":2,"node_type":"ms"}]`, azName)

	updatedSpecAzInfo := fmt.Sprintf(`[{"availability_zone_name":"%s","availability_zone_count":3,"node_type":"ms"}]`, azName)
	//backupStorageType := "SATA"
	//updatedBackupStorageSpace := 160
	//updatedAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"s"}]`

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
			// 创建一个单节点的mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort, storageType, storageSpace,
					backupStorageType, azInfo),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			// 更新mongodb实例，升级存储空间和flavor_name
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, prodId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedSpecAzInfo),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(updatedReadPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			// 更新mongodb实例，升级存储空间
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, updatedProdId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedAzInfo),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(updatedReadPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, updatedProdId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedAzInfo),
				Destroy: true,
			},
		},
	})
}

// 副本集，不传azList, 存储为SATA
func TestAccCtyunMongodbInstanceReplicaSATANoAzList(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	//dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mongodb_instance." + rnd
	//datasourceName := "data.ctyun_mongodb_instances." + dnd

	resourceFile := "resource_ctyun_mongodb_instance_replica_on_demand_no_az.tf"
	// datasourceFile := "datasource_ctyun_mongodb_instances.tf"
	// 创建参数
	cycleType := "on_demand"
	vpcID := dependence.vpcID
	flavorName := "s7.large.2"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mongodb-single-" + utils.GenerateRandomString()
	password := "Kyk123=" + utils.GenerateRandomString()
	prodId := "Replica3R34"
	readPort := 12345
	storageType := "SAS"
	storageSpace := 100
	backupStorageType := "SATA"
	//backupStorageSpace := 120
	//azInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"master"},
	//			{"availability_zone_name":"cn-huadong1-jsnj2A-public-ctcloud","availability_zone_count":1,"node_type":"master"},
	//			{"availability_zone_name":"cn-huadong1-jsnj3A-public-ctcloud","availability_zone_count":1,"node_type":"master"}]`

	//更新参数
	updatedName := "tf-mongodb-single-new-" + utils.GenerateRandomString()
	updatedFlavorName := "s7.large.4"
	updatedReadPort := 12348
	//updatedStorageType := ""
	updatedStorageSpace := 110
	//updatedAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":2,"node_type":"master"}]`
	//backupStorageType := "SATA"
	//updatedBackupStorageSpace := 160
	//updatedAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"s"}]`

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
			// 创建一个单节点的mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort, storageType, storageSpace,
					backupStorageType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			// 更新mongodb实例，升级存储空间
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, prodId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(updatedReadPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, prodId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType),
				Destroy: true,
			},
		},
	})
}

// 集群版，传azList，更新端口、名称和主存储空间
func TestAccCtyunMongodbInstanceClusterOs(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	//dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mongodb_instance." + rnd
	//datasourceName := "data.ctyun_mongodb_instances." + dnd

	resourceFile := "resource_ctyun_mongodb_instance_cluster_on_demand_os.tf"
	//resourceFile1 := "resource_ctyun_mongodb_instance_cluster_on_demand_os_update.tf"
	//datasourceFile := "datasource_ctyun_mongodb_instances.tf"
	// 创建参数
	cycleType := "on_demand"
	vpcID := dependence.vpcID
	flavorName := "s7.large.2"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mongodb-single-" + utils.GenerateRandomString()
	password := "Kyk123=" + utils.GenerateRandomString()
	prodId := "Cluster34"
	readPort := 12345
	storageType := "SAS"
	storageSpace := 100
	backupStorageType := "OS"
	azName := dependence.azName
	azInfo := fmt.Sprintf(`[{"availability_zone_name":"%s","availability_zone_count":2,"node_type":"mongos"},
				{"availability_zone_name":"%s","availability_zone_count":6,"node_type":"shard"},
				{"availability_zone_name":"%s","availability_zone_count":3,"node_type":"config"}]`, azName, azName, azName)
	shardNum := 2
	mongosNum := 2

	//更新参数
	updatedName := "tf-mongodb-single-new-" + utils.GenerateRandomString()
	updatedReadPort := 12348

	//updatedStorageType := ""
	updatedStorageSpace := 110
	//updatedSpecAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":2,"node_type":"ms"}]`

	//backupStorageType := "SATA"
	//updatedBackupStorageSpace := 160
	//updatedAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"s"}]`

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
			// 创建一个单节点的mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort, storageType, storageSpace,
					backupStorageType, azInfo, shardNum, mongosNum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "shard_num", strconv.Itoa(shardNum)),
					resource.TestCheckResourceAttr(resourceName, "mongos_num", strconv.Itoa(mongosNum)),
				),
			},
			// 更新mongodb实例，升级存储空间、 name 和port
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, updatedName, password, prodId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, azInfo, shardNum, mongosNum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(updatedReadPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, updatedName, password, prodId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, azInfo, shardNum, mongosNum),
				Destroy: true,
			},
		},
	})
}

// 集群版，传azList，更新mongos 和shard spec
func TestAccCtyunMongodbInstanceClusterOsUpdateMongosSpec(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	//dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mongodb_instance." + rnd
	//datasourceName := "data.ctyun_mongodb_instances." + dnd

	resourceFile := "resource_ctyun_mongodb_instance_cluster_on_demand_os.tf"
	resourceFile1 := "resource_ctyun_mongodb_instance_cluster_on_demand_os_update.tf"
	//datasourceFile := "datasource_ctyun_mongodb_instances.tf"
	// 创建参数
	cycleType := "on_demand"
	vpcID := dependence.vpcID
	flavorName := "s7.large.2"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mongodb-single-" + utils.GenerateRandomString()
	password := "Kyk123=" + utils.GenerateRandomString()
	prodId := "Cluster34"
	readPort := 12345
	storageType := "SAS"
	storageSpace := 100
	backupStorageType := "OS"
	azName := dependence.azName
	azInfo := fmt.Sprintf(`[{"availability_zone_name":"%s","availability_zone_count":2,"node_type":"mongos"},
				{"availability_zone_name":"%s","availability_zone_count":6,"node_type":"shard"},
				{"availability_zone_name":"%s","availability_zone_count":3,"node_type":"config"}]`, azName, azName, azName)
	shardNum := 2
	mongosNum := 2

	//更新参数
	updatedFlavorName := "s7.large.4"

	updatedMongosSpecAzInfo := fmt.Sprintf(`[{"availability_zone_name":"%s","availability_zone_count":2,"node_type":"mongos"}]`, azName)
	updatedShardSpecAzInfo := fmt.Sprintf(`[{"availability_zone_name":"%s","availability_zone_count":6,"node_type":"shard"}]`, azName)

	upgradeNodeTypeMongos := "mongos"
	upgradeNodeTypeShard := "shard"

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
			// 创建一个单节点的mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort, storageType, storageSpace,
					backupStorageType, azInfo, shardNum, mongosNum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "shard_num", strconv.Itoa(shardNum)),
					resource.TestCheckResourceAttr(resourceName, "mongos_num", strconv.Itoa(mongosNum)),
				),
			},
			// 更新mongodb实例，升级存储空间、mongos的spec
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, name, password, prodId, readPort,
					storageType, storageSpace, backupStorageType, updatedMongosSpecAzInfo, shardNum, mongosNum, upgradeNodeTypeMongos),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			// 扩容
			// 更新mongodb实例，shard的spec
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, name, password, prodId, readPort,
					storageType, storageSpace, backupStorageType, updatedShardSpecAzInfo, shardNum, mongosNum, upgradeNodeTypeShard),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(30 * time.Second)
							return nil
						},
					),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, name, password, prodId, readPort,
					storageType, storageSpace, backupStorageType, updatedShardSpecAzInfo, shardNum, mongosNum, upgradeNodeTypeShard),
				Destroy: true,
			},
		},
	})
}

// 集群版，传azList，更新mongos shard 节点数
func TestAccCtyunMongodbInstanceClusterOsUpdateNodeNum(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	//dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mongodb_instance." + rnd
	//datasourceName := "data.ctyun_mongodb_instances." + dnd

	resourceFile := "resource_ctyun_mongodb_instance_cluster_on_demand_os.tf"
	resourceFile1 := "resource_ctyun_mongodb_instance_cluster_on_demand_os_update.tf"
	//datasourceFile := "datasource_ctyun_mongodb_instances.tf"
	// 创建参数
	cycleType := "on_demand"
	vpcID := dependence.vpcID
	flavorName := "s7.large.2"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mongodb-single-" + utils.GenerateRandomString()
	password := "Kyk123=" + utils.GenerateRandomString()
	prodId := "Cluster34"
	readPort := 12345
	storageType := "SAS"
	storageSpace := 100
	backupStorageType := "OS"
	azName := dependence.azName
	azInfo := fmt.Sprintf(`[{"availability_zone_name":"%s","availability_zone_count":2,"node_type":"mongos"},
				{"availability_zone_name":"%s","availability_zone_count":6,"node_type":"shard"},
				{"availability_zone_name":"%s","availability_zone_count":3,"node_type":"config"}]`, azName, azName, azName)
	shardNum := 2
	mongosNum := 2

	//更新参数

	updatedMongosNodeAzInfo := fmt.Sprintf(`[{"availability_zone_name":"%s","availability_zone_count":1,"node_type":"mongos"}]`, azName)
	updatedShardNodeAzInfo := fmt.Sprintf(`[{"availability_zone_name":"%s","availability_zone_count":3,"node_type":"shard"}]`, azName)
	upgradeNodeTypeShard := "shard"
	upgradeNodeTypeMongos := "mongos"

	updatedShardNum := 3
	updatedMongosNum := 3

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
			// 创建一个单节点的mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort, storageType, storageSpace,
					backupStorageType, azInfo, shardNum, mongosNum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "shard_num", strconv.Itoa(shardNum)),
					resource.TestCheckResourceAttr(resourceName, "mongos_num", strconv.Itoa(mongosNum)),
				),
			},
			// 更新mongodb实例，升级存储空间、mongos的节点
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort,
					storageType, storageSpace, backupStorageType, updatedMongosNodeAzInfo, shardNum, updatedMongosNum, upgradeNodeTypeMongos),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			// 扩容
			// 更新mongodb实例，shard的节点
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort,
					storageType, storageSpace, backupStorageType, updatedShardNodeAzInfo, updatedShardNum, updatedMongosNum, upgradeNodeTypeShard),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort,
					storageType, storageSpace, backupStorageType, updatedShardNodeAzInfo, updatedShardNum, updatedMongosNum, upgradeNodeTypeShard),
				Destroy: true,
			},
		},
	})
}

// 集群版，不传azList,修改存储，备份空间，端口
func TestAccCtyunMongodbInstanceClusterNoAz(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	//dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mongodb_instance." + rnd
	//datasourceName := "data.ctyun_mongodb_instances." + dnd

	resourceFile := "resource_ctyun_mongodb_instance_cluster_on_demand_no_az.tf"
	//resourceFile1 := "resource_ctyun_mongodb_instance_cluster_on_demand_no_az_update.tf"
	//datasourceFile := "datasource_ctyun_mongodb_instances.tf"
	// 创建参数
	cycleType := "on_demand"
	vpcID := dependence.vpcID
	flavorName := "s7.large.2"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mongodb-single-" + utils.GenerateRandomString()
	password := "Kyk123=" + utils.GenerateRandomString()
	prodId := "Cluster40"
	readPort := 12345
	storageType := "SAS"
	storageSpace := 100
	backupStorageType := "SATA"
	//backupStorageSpace := 120
	shardNum := 2
	mongosNum := 2

	//更新参数
	updatedName := "tf-mongodb-single-new-" + utils.GenerateRandomString()
	updatedReadPort := 12348

	//updatedStorageType := ""
	updatedStorageSpace := 110
	//updatedSpecAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":2,"node_type":"ms"}]`
	//updatedProdIDAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":2,"node_type":"ms"}]`
	//
	//backupStorageType := "SATA"
	//updatedBackupStorageSpace := 150
	//updatedAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"s"}]`

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
			// 创建一个单节点的mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort, storageType, storageSpace,
					backupStorageType, shardNum, mongosNum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "shard_num", strconv.Itoa(shardNum)),
					resource.TestCheckResourceAttr(resourceName, "mongos_num", strconv.Itoa(mongosNum)),
				),
			},
			// 更新mongodb实例，升级存储空间、mongos的spec
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, updatedName, password, prodId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, shardNum, mongosNum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(updatedReadPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, updatedName, password, prodId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, shardNum, mongosNum),
				Destroy: true,
			},
		},
	})
}

// mongodb升配mongos spec
func TestAccCtyunMongodbInstanceClusterNoAzUpdateMongosSpec(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	//dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mongodb_instance." + rnd
	//datasourceName := "data.ctyun_mongodb_instances." + dnd

	resourceFile := "resource_ctyun_mongodb_instance_cluster_on_demand_no_az.tf"
	resourceFile1 := "resource_ctyun_mongodb_instance_cluster_on_demand_no_az_update.tf"
	//datasourceFile := "datasource_ctyun_mongodb_instances.tf"
	// 创建参数
	cycleType := "on_demand"
	vpcID := dependence.vpcID
	flavorName := "s7.large.2"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mongodb-single-" + utils.GenerateRandomString()
	password := "Kyk123=" + utils.GenerateRandomString()
	prodId := "Cluster40"
	readPort := 12345
	storageType := "SAS"
	storageSpace := 100
	backupStorageType := "SATA"
	//backupStorageSpace := 120
	shardNum := 2
	mongosNum := 2

	updatedFlavorName := "s7.large.4"

	//updatedStorageType := ""
	//updatedSpecAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":2,"node_type":"ms"}]`
	//updatedProdIDAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":2,"node_type":"ms"}]`
	//
	upgradeNodeTypeMongos := "mongos"

	//backupStorageType := "SATA"
	//updatedAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"s"}]`

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
			// 创建一个单节点的mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort, storageType, storageSpace,
					backupStorageType, shardNum, mongosNum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "shard_num", strconv.Itoa(shardNum)),
					resource.TestCheckResourceAttr(resourceName, "mongos_num", strconv.Itoa(mongosNum)),
				),
			},
			// 更新mongodb实例、mongos的spec升配
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, name, password, prodId, readPort,
					storageType, storageSpace, backupStorageType, shardNum, mongosNum, upgradeNodeTypeMongos),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			// 扩容
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, name, password, prodId, readPort,
					storageType, storageSpace, backupStorageType, shardNum, mongosNum, upgradeNodeTypeMongos),
				Destroy: true,
			},
		},
	})
}

// mongodb升配 shard spec
func TestAccCtyunMongodbInstanceClusterNoAzUpdateShardSpec(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	//dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mongodb_instance." + rnd
	//datasourceName := "data.ctyun_mongodb_instances." + dnd

	resourceFile := "resource_ctyun_mongodb_instance_cluster_on_demand_no_az.tf"
	resourceFile1 := "resource_ctyun_mongodb_instance_cluster_on_demand_no_az_update.tf"
	//datasourceFile := "datasource_ctyun_mongodb_instances.tf"
	// 创建参数
	cycleType := "on_demand"
	vpcID := dependence.vpcID
	flavorName := "s7.large.2"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mongodb-single-" + utils.GenerateRandomString()
	password := "Kyk123=" + utils.GenerateRandomString()
	prodId := "Cluster40"
	readPort := 12345
	storageType := "SAS"
	storageSpace := 100
	backupStorageType := "SATA"
	//backupStorageSpace := 120
	shardNum := 2
	mongosNum := 2

	updatedFlavorName := "s7.large.4"

	//updatedStorageType := ""
	//updatedSpecAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":2,"node_type":"ms"}]`
	//updatedProdIDAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":2,"node_type":"ms"}]`
	//
	upgradeNodeTypeShard := "shard"

	//backupStorageType := "SATA"
	//updatedAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"s"}]`

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
			// 创建一个单节点的mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort, storageType, storageSpace,
					backupStorageType, shardNum, mongosNum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "shard_num", strconv.Itoa(shardNum)),
					resource.TestCheckResourceAttr(resourceName, "mongos_num", strconv.Itoa(mongosNum)),
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(30 * time.Second)
							return nil
						},
					),
				),
			},
			// 扩容
			// 更新mongodb实例，shard的spec
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, name, password, prodId, readPort,
					storageType, storageSpace, backupStorageType, shardNum, mongosNum, upgradeNodeTypeShard),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(30 * time.Second)
							return nil
						},
					),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, name, password, prodId, readPort,
					storageType, storageSpace, backupStorageType, shardNum, mongosNum, upgradeNodeTypeShard),
				Destroy: true,
			},
		},
	})
}

// 集群版，不传azList,升配shard和mongos节点数量
func TestAccCtyunMongodbInstanceClusterNoAzUpdateNode(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	//dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mongodb_instance." + rnd
	//datasourceName := "data.ctyun_mongodb_instances." + dnd

	resourceFile := "resource_ctyun_mongodb_instance_cluster_on_demand_no_az.tf"
	resourceFile1 := "resource_ctyun_mongodb_instance_cluster_on_demand_no_az_update.tf"
	//datasourceFile := "datasource_ctyun_mongodb_instances.tf"
	// 创建参数
	cycleType := "on_demand"
	vpcID := dependence.vpcID
	flavorName := "s7.large.2"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mongodb-single-" + utils.GenerateRandomString()
	password := "Kyk123=" + utils.GenerateRandomString()
	prodId := "Cluster40"
	readPort := 12345
	storageType := "SAS"
	storageSpace := 100
	backupStorageType := "SATA"
	shardNum := 2
	mongosNum := 2

	//更新参数

	//updatedStorageType := ""
	//updatedSpecAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":2,"node_type":"ms"}]`
	//updatedProdIDAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":2,"node_type":"ms"}]`
	updatedShardNum := 3
	updatedMongosNum := 3
	//
	upgradeNodeTypeMongos := "mongos"
	upgradeNodeTypeShard := "shard"

	//backupStorageType := "SATA"
	//updatedAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"s"}]`

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
			// 创建一个单节点的mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort, storageType, storageSpace,
					backupStorageType, shardNum, mongosNum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "shard_num", strconv.Itoa(shardNum)),
					resource.TestCheckResourceAttr(resourceName, "mongos_num", strconv.Itoa(mongosNum)),
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(30 * time.Second)
							return nil
						},
					),
				),
			},
			// 扩容
			// 更新shard数量
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort,
					storageType, storageSpace, backupStorageType, updatedShardNum, mongosNum, upgradeNodeTypeShard),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "shard_num", strconv.Itoa(updatedShardNum)),
					resource.TestCheckResourceAttr(resourceName, "mongos_num", strconv.Itoa(mongosNum)),
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(30 * time.Second)
							return nil
						},
					),
				),
			},
			// 更新mongos数量
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort,
					storageType, storageSpace, backupStorageType, updatedShardNum, updatedMongosNum, upgradeNodeTypeMongos),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "shard_num", strconv.Itoa(updatedShardNum)),
					resource.TestCheckResourceAttr(resourceName, "mongos_num", strconv.Itoa(updatedMongosNum)),
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(30 * time.Second)
							return nil
						},
					),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort,
					storageType, storageSpace, backupStorageType, updatedShardNum, updatedMongosNum, upgradeNodeTypeMongos),
				Destroy: true,
			},
		},
	})
}
