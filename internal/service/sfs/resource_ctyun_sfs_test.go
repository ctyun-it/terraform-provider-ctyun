package sfs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunSfs(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_sfs." + rnd
	resourceFile := "resource_ctyun_sfs_onDemand.tf"
	resourceFile1 := "resource_ctyun_sfs_onDemand_readonly.tf"

	// 配置测试环境需要的动态值（实际测试时替换为有效值）
	azName := "cn-huadong1-jsnj1A-public-ctcloud"
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	sfsType := "performance"
	sfsProtocol := "nfs"
	name := "sfs-" + utils.GenerateRandomString()
	sfsSize := 500
	cycleType := "on_demand"

	updatedSfsSize := 550
	updatedName := name + "new"
	readOnly := true

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
			// 1. 基础创建测试
			{
				Config: utils.LoadTestCase(resourceFile, rnd, sfsType, sfsProtocol, name, sfsSize, cycleType, azName, vpcID, subnetID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "sfs_size", fmt.Sprintf("%d", sfsSize)),
					resource.TestCheckResourceAttr(resourceName, "sfs_protocol", sfsProtocol),
					resource.TestCheckResourceAttr(resourceName, "sfs_type", sfsType),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "read_only", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			// 2. 资源更新测试（名称/大小/只读）
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, sfsType, sfsProtocol, updatedName, updatedSfsSize, cycleType, azName, vpcID, subnetID, readOnly),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "sfs_size", fmt.Sprintf("%d", updatedSfsSize)),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "read_only", "true"),
				),
			},
			// 3. 资源导入测试
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, sfsType, sfsProtocol, updatedName, updatedSfsSize, cycleType, azName, vpcID, subnetID, readOnly),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunSfsCycle(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_sfs." + rnd
	resourceFile := "resource_ctyun_sfs_cycle.tf"

	datasourceName := "data.ctyun_sfs_instances." + dnd
	datasourceFile := "datasource_ctyun_sfs_instances.tf"

	// 配置测试环境需要的动态值
	//azName := "cn-huadong1-jsnj1A-public-ctcloud"
	vpcID := dependence.vpcID
	subnetID := dependence.subnetID
	//kmsUUID := "e8c81488-7990-4123-b79e-1235d2b1f4eb" // 添加KMS UUID

	// 加密相关参数
	//isEncrypt := true
	isEncrypt := false

	// 存储类型和协议
	sfsType := "performance"
	sfsProtocol := "cifs"

	// 命名和大小
	name := "sfs-" + rnd
	sfsSize := 500
	cycleType := "month"
	cycleCount := 1
	//readOnly := true

	// 更新参数
	updatedSfsSize := 550
	updatedName := name + "-updated"
	//updatedReadOnly := false

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
			// 1. 基础创建测试（加密性能型CIFS存储）
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd, isEncrypt, sfsType, sfsProtocol,
					name, sfsSize, cycleType, cycleCount, vpcID, subnetID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "sfs_size", fmt.Sprintf("%d", sfsSize)),
					resource.TestCheckResourceAttr(resourceName, "sfs_protocol", sfsProtocol),
					resource.TestCheckResourceAttr(resourceName, "sfs_type", sfsType),
					resource.TestCheckResourceAttr(resourceName, "cycle_type", cycleType),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "is_encrypt", "false"), // 验证加密
					resource.TestCheckResourceAttr(resourceName, "read_only", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			// 2. 资源更新测试（名称/大小/只读）
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd, isEncrypt, sfsType, sfsProtocol,
					updatedName, updatedSfsSize, cycleType, cycleCount, vpcID, subnetID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "sfs_size", fmt.Sprintf("%d", updatedSfsSize)),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "read_only", "false"),
					// 保持加密设置不变
					resource.TestCheckResourceAttr(resourceName, "is_encrypt", "false"),
				),
			},

			// 3. datasource 验证
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd, isEncrypt, sfsType, sfsProtocol,
					updatedName, updatedSfsSize, cycleType, cycleCount, vpcID, subnetID) +
					utils.LoadTestCase(datasourceFile, dnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "region_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "sfs_list.#"),
					//resource.TestCheckResourceAttr(datasourceName, "sfs_list.0.sfs_name", updatedName),
					//resource.TestCheckResourceAttr(datasourceName, "sfs_list.0.sfs_size", fmt.Sprintf("%d", updatedSfsSize)),
					//resource.TestCheckResourceAttr(datasourceName, "sfs_list.0.sfs_type", sfsType),
					//resource.TestCheckResourceAttr(datasourceName, "sfs_list.0.sfs_protocol", sfsProtocol),
					//resource.TestCheckResourceAttr(datasourceName, "sfs_list.0.is_encrypt", "false"),
				),
			},
			// 4. 资源导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					projectId := ds.Attributes["project_id"]
					if id == "" || regionId == "" {
						return "", fmt.Errorf("id or region_id is required")
					}
					return fmt.Sprintf("%s,%s,%s", id, regionId, projectId), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cycle_count", "kms_uuid", "cycle_type", "is_encrypt", "vpc_id", "subnet_id", "az_name", "used_size"},
			},
			// 4. 清理资源
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd, isEncrypt, sfsType, sfsProtocol,
					updatedName, updatedSfsSize, cycleType, cycleCount, vpcID, subnetID),
				Destroy: true,
			},
		},
	})
}
