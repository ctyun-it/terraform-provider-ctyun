package opensearch_test

import (
	"fmt"
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunSearchInstance_OpenSearch(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_search_instance." + rnd
	datasourceName := "data.ctyun_search_instances." + rnd
	resourceFile := "resource_ctyun_search_instance.tf"
	datasourceFile := "datasource_ctyun_search_instances.tf"

	clusterName := utils.GenerateRandomString()
	regionID := dependence.regionID
	zoneList := dependence.azName
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	enableIPv6 := true
	clusterType := 1
	osType := "CTyun"
	enableHTTPS := "CLOSE"
	hostNum := 3
	ioType := dependence.storageType
	volume := 40
	iaasVmSpecCode := dependence.flavorName
	cycleCnt := 1
	cycleType := "month"
	nodeGroupType := "MASTER"
	exclusiveMaster := "EXCLUSIVE_MASTER"
	coordinate := "COORDINATE"
	cold := "COLD"

	// 生成随机密码
	randomPwd := "Th@5" + utils.GenerateRandomString()
	// 确保密码长度在 12-26 位之间
	componentPwd := randomPwd[:min(len(randomPwd), 20)]

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
				// Step 1: 创建实例
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					clusterName,
					regionID,
					zoneList,
					vpcID,
					subnetID,
					securityGroupID,
					enableIPv6,
					componentPwd,
					clusterType,
					osType,
					enableHTTPS,
					cycleCnt,
					cycleType,
					hostNum,
					ioType,
					volume,
					iaasVmSpecCode,
					nodeGroupType,
					exclusiveMaster,
					coordinate,
					cold,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				// Step 2: 更新节点配置（扩容）
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					clusterName,
					regionID,
					zoneList,
					vpcID,
					subnetID,
					securityGroupID,
					enableIPv6,
					componentPwd,
					clusterType,
					osType,
					enableHTTPS,
					cycleCnt,
					cycleType,
					4,
					ioType,
					volume,
					iaasVmSpecCode,
					nodeGroupType,
					exclusiveMaster,
					coordinate,
					cold,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
				),
			},
			{
				// Step 3: 导入测试
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "auto_pay", "cycle_type", "enable_https"},
			},
			{
				// Step 4: Datasource 数据源测试
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					clusterName,
					regionID,
					zoneList,
					vpcID,
					subnetID,
					securityGroupID,
					enableIPv6,
					componentPwd,
					clusterType,
					osType,
					enableHTTPS,
					cycleCnt,
					cycleType,
					4,
					ioType,
					volume,
					iaasVmSpecCode,
					nodeGroupType,
					exclusiveMaster,
					coordinate,
					cold,
				) + utils.LoadTestCase(datasourceFile, rnd, regionID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "total"),
				),
			},
			{
				// Step 2: 更新节点配置（扩容）
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					clusterName,
					regionID,
					zoneList,
					vpcID,
					subnetID,
					securityGroupID,
					enableIPv6,
					componentPwd,
					clusterType,
					osType,
					enableHTTPS,
					cycleCnt,
					cycleType,
					4,
					ioType,
					volume,
					iaasVmSpecCode,
					nodeGroupType,
					exclusiveMaster,
					coordinate,
					cold,
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunSearchInstance_ElasticSearch(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_search_instance." + rnd
	datasourceName := "data.ctyun_search_instances." + rnd
	resourceFile := "resource_ctyun_search_instance.tf"
	datasourceFile := "datasource_ctyun_search_instances.tf"

	clusterName := utils.GenerateRandomString()
	zoneList := dependence.azName
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	// 生成随机密码
	randomPwd := "Th@5" + utils.GenerateRandomString()
	// 确保密码长度在 12-26 位之间
	componentPwd := randomPwd[:min(len(randomPwd), 20)]
	regionID := dependence.regionID
	enableIPv6 := false
	clusterType := 2
	osType := "CTyun"
	enableHTTPS := "CLOSE"
	cycleCnt := 1
	cycleType := "month"
	hostNum := 3
	ioType := dependence.storageType
	volume := 40
	iaasVmSpecCode := dependence.flavorName
	nodeGroupType := "MASTER"
	exclusiveMaster := "EXCLUSIVE_MASTER"
	coordinate := "COORDINATE"
	cold := "COLD"

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
				// Step 1: 创建实例
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					clusterName,
					regionID,
					zoneList,
					vpcID,
					subnetID,
					securityGroupID,
					enableIPv6,
					componentPwd,
					clusterType,
					osType,
					enableHTTPS,
					cycleCnt,
					cycleType,
					hostNum,
					ioType,
					volume,
					iaasVmSpecCode,
					nodeGroupType,
					exclusiveMaster,
					coordinate,
					cold,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				// Step 2: 更新节点配置（扩容）
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					clusterName,
					regionID,
					zoneList,
					vpcID,
					subnetID,
					securityGroupID,
					enableIPv6,
					componentPwd,
					clusterType,
					osType,
					enableHTTPS,
					cycleCnt,
					cycleType,
					4,
					ioType,
					volume,
					iaasVmSpecCode,
					nodeGroupType,
					exclusiveMaster,
					coordinate,
					cold,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
				),
			},
			{
				// Step 3: 导入测试
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "auto_pay", "cycle_cnt", "cycle_type", "enable_https", "pay_type"},
			},
			{
				// Step 4: Datasource 数据源测试
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					clusterName,
					regionID,
					zoneList,
					vpcID,
					subnetID,
					securityGroupID,
					enableIPv6,
					componentPwd,
					clusterType,
					osType,
					enableHTTPS,
					cycleCnt,
					cycleType,
					4,
					ioType,
					volume,
					iaasVmSpecCode,
					nodeGroupType,
					exclusiveMaster,
					coordinate,
					cold,
				) + utils.LoadTestCase(datasourceFile, rnd, regionID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "total"),
				),
			},
			{
				// Step 5: 退订
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					clusterName,
					regionID,
					zoneList,
					vpcID,
					subnetID,
					securityGroupID,
					enableIPv6,
					componentPwd,
					clusterType,
					osType,
					enableHTTPS,
					cycleCnt,
					cycleType,
					4,
					ioType,
					volume,
					iaasVmSpecCode,
					nodeGroupType,
					exclusiveMaster,
					coordinate,
					cold,
				),
				Destroy: true,
			},
		},
	})
}
