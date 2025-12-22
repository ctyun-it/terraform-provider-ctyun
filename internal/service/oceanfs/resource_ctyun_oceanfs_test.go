package oceanfs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunOceanfs(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_oceanfs." + rnd
	resourceFile := "resource_ctyun_oceanfs_no_tags.tf"

	// 配置测试环境需要的动态值（实际测试时替换为有效值）
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	//sfsProtocol := "nfs"
	sfsSize := 100 // 最小值100GB
	cycleType := "on_demand"
	projectID := "0"

	updatedSfsSize := 150
	proctcols := []string{"nfs", "cifs"}
	for _, sfsProtocol := range proctcols {
		fmt.Println(fmt.Sprintf("protocol= %s 的oceanfs 验证", sfsProtocol))
		name := "oceanfs-" + utils.GenerateRandomString() + "-" + sfsProtocol
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
				// 1. 基础创建测试（按需计费NFS）
				{
					Config: utils.LoadTestCase(resourceFile, rnd, projectID, sfsProtocol, name, sfsSize, cycleType, vpcID, subnetID),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttr(resourceName, "size", fmt.Sprintf("%d", sfsSize)),
						resource.TestCheckResourceAttr(resourceName, "protocol", sfsProtocol),
						resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
						resource.TestCheckResourceAttr(resourceName, "name", name),
						resource.TestCheckResourceAttrSet(resourceName, "status"),
						resource.TestCheckResourceAttrSet(resourceName, "create_time"),
						resource.TestCheckResourceAttrSet(resourceName, "used_size"),
					),
				},
				// 2. 资源更新测试（更新大小和名称）
				{
					Config: utils.LoadTestCase(resourceFile, rnd, projectID, sfsProtocol, name, updatedSfsSize, cycleType, vpcID, subnetID),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestCheckResourceAttr(resourceName, "size", fmt.Sprintf("%d", updatedSfsSize)),
						resource.TestCheckResourceAttr(resourceName, "name", name),
					),
				},
				// 3. import state验证 1
				{
					ResourceName: resourceName,
					ImportState:  true,
					ImportStateIdFunc: func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceName)
						}
						return fmt.Sprintf("%s,%s,%s",
							rs.Primary.Attributes["id"],
							rs.Primary.Attributes["project_id"],
							rs.Primary.Attributes["region_id"],
						), nil
					},
					ImportStateVerify: true,
					ImportStateVerifyIgnore: []string{
						"tags",
						"cycle_type",
						"subnet_id",
						"vpc_id",
						"update_time",
					},
				},
				// 4. import state验证 2
				{
					ResourceName: resourceName,
					ImportState:  true,
					ImportStateIdFunc: func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceName)
						}
						return fmt.Sprintf("%s",
							rs.Primary.Attributes["id"],
						), nil
					},
					ImportStateVerify: true,
					ImportStateVerifyIgnore: []string{
						"tags",
						"cycle_type",
						"subnet_id",
						"vpc_id",
						"update_time",
						"project_id",
					},
				},
				// 5.销毁资源
				{
					Config:  utils.LoadTestCase(resourceFile, rnd, projectID, sfsProtocol, name, updatedSfsSize, cycleType, vpcID, subnetID),
					Destroy: true,
				},
			},
		})
	}
}

func TestAccCtyunOceanfsWithVpce(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_oceanfs." + rnd
	resourceFile := "resource_ctyun_oceanfs_vpce.tf"
	datasourceFile := "resource_ctyun_oceanfs_instances.tf"
	datasourceName := "data.ctyun_oceanfs_instances." + dnd

	// 配置测试环境需要的动态值
	//azName := "cn-huabei2-tj1A-public-ctcloud" // 替换为实际可用区
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	sfsProtocol := "nfs"
	name := "oceanfs-vpce-" + utils.GenerateRandomString()
	sfsSize := 100
	cycleType := "on_demand"
	isVpce := true

	updatedSfsSize := 120

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
			// 1. 创建带VPCE的OceanFS
			{
				Config: utils.LoadTestCase(resourceFile, rnd, sfsProtocol, name, sfsSize, cycleType, vpcID, subnetID, isVpce),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "size", fmt.Sprintf("%d", sfsSize)),
					resource.TestCheckResourceAttr(resourceName, "protocol", sfsProtocol),
					resource.TestCheckResourceAttr(resourceName, "is_vpce", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			// 2. update
			{
				Config: utils.LoadTestCase(resourceFile, rnd, sfsProtocol, name, updatedSfsSize, cycleType, vpcID, subnetID, isVpce),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "size", fmt.Sprintf("%d", updatedSfsSize)),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			// 3. datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, sfsProtocol, name, updatedSfsSize, cycleType, vpcID, subnetID, isVpce) +
					utils.LoadTestCase(datasourceFile, dnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "instances.#"),
				),
			},
			// 4. 销毁资源
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, sfsProtocol, name, sfsSize, cycleType, vpcID, subnetID, isVpce),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunOceanfsCycle(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_oceanfs." + rnd
	resourceFile := "resource_ctyun_oceanfs_period.tf"

	// 配置测试环境需要的动态值
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	sfsProtocol := "nfs"
	name := "oceanfs-cycle-" + utils.GenerateRandomString()
	sfsSize := 100
	cycleType := "month"
	cycleCount := 1
	tags := fmt.Sprintf(`{"key":"test","value":"%s"},{"key":"test1","value":"%s"}`, rnd, rnd)

	updatedSfsSize := 120

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
			// 1. 基础创建测试（包月计费）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, sfsProtocol, name, sfsSize, cycleType, cycleCount, vpcID, subnetID, tags),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "size", fmt.Sprintf("%d", sfsSize)),
					resource.TestCheckResourceAttr(resourceName, "protocol", sfsProtocol),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "cycle_count", fmt.Sprintf("%d", cycleCount)),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "expire_time"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
				),
			},
			// 2. 资源更新测试（仅更新大小）
			{
				Config: utils.LoadTestCase(resourceFile, rnd, sfsProtocol, name, updatedSfsSize, cycleType, cycleCount, vpcID, subnetID, tags),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "size", fmt.Sprintf("%d", updatedSfsSize)),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, sfsProtocol, name, updatedSfsSize, cycleType, cycleCount, vpcID, subnetID, tags),
				Destroy: true,
			},
		},
	})
}
