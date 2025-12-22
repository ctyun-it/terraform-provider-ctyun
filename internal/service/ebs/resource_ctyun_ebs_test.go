package ebs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunEbs(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	and := utils.GenerateRandomString()

	resourceName := "ctyun_ebs." + rnd
	datasourceName := "data.ctyun_ebs_volumes." + dnd
	resourceFile := "resource_ctyun_ebs.tf"
	datasourceFile := "datasource_ctyun_ebs_volumes.tf"
	associationFile := "resource_ctyun_ebs_association_ecs.tf"

	associationResourceName := "ctyun_ebs_association_ecs." + and
	initName := "init-ebs-" + rnd
	initSize := 60

	updatedName := "updated-ebs-" + rnd
	updatedSize := 100

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

			// 创建
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initName,
					"sata",
					initSize,
					"",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "size", strconv.Itoa(initSize)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			// 更新属性，同时绑定ecs
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					"sata",
					updatedSize,
					"",
				) + utils.LoadTestCase(
					associationFile, and,
					dependence.ecsID,
					resourceName+".id",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "size", strconv.Itoa(updatedSize)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},

			// 通过查询检查是否绑定成功
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					"sata",
					updatedSize,
					"",
				) + utils.LoadTestCase(
					associationFile, and,
					dependence.ecsID,
					resourceName+".id",
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					"",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					func(s *terraform.State) error {
						ds := s.RootModule().Resources[datasourceName].Primary

						count, err := strconv.Atoi(ds.Attributes["volumes.#"])
						if err != nil || count == 0 {
							return fmt.Errorf("volumes 无效: %v", ds.Attributes)
						}

						for i := 0; i < count; i++ {
							if ds.Attributes[fmt.Sprintf("volumes.%d.name", i)] == updatedName {
								if dependence.ecsID == ds.Attributes[fmt.Sprintf("volumes.%d.attachments.0.instance_id", i)] {
									return nil
								} else {
									return fmt.Errorf("绑定云主机失败")
								}
							}
						}
						return fmt.Errorf("未找到目标元素")
					},
				),
			},
			// 通过查询检查是否绑定成功
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					"sata",
					updatedSize,
					"",
				) + utils.LoadTestCase(
					associationFile, and,
					dependence.ecsID,
					resourceName+".id",
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					fmt.Sprintf("disk_id = \"%s\"\n", "${"+resourceName+".id}"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					func(s *terraform.State) error {
						ds := s.RootModule().Resources[datasourceName].Primary

						count, err := strconv.Atoi(ds.Attributes["volumes.#"])
						if err != nil || count == 0 {
							return fmt.Errorf("volumes 无效: %v", ds.Attributes)
						}

						for i := 0; i < count; i++ {
							if ds.Attributes[fmt.Sprintf("volumes.%d.name", i)] == updatedName {
								if dependence.ecsID == ds.Attributes[fmt.Sprintf("volumes.%d.attachments.0.instance_id", i)] {
									return nil
								} else {
									return fmt.Errorf("绑定云主机失败")
								}
							}
						}
						return fmt.Errorf("未找到目标元素")
					},
				),
			},
			// 通过查询检查是否绑定成功
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					"sata",
					updatedSize,
					"",
				) + utils.LoadTestCase(
					associationFile, and,
					dependence.ecsID,
					resourceName+".id",
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					fmt.Sprintf("disk_name = \"%s\"\n", updatedName),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					func(s *terraform.State) error {
						ds := s.RootModule().Resources[datasourceName].Primary

						count, err := strconv.Atoi(ds.Attributes["volumes.#"])
						if err != nil || count == 0 {
							return fmt.Errorf("volumes 无效: %v", ds.Attributes)
						}

						for i := 0; i < count; i++ {
							if ds.Attributes[fmt.Sprintf("volumes.%d.name", i)] == updatedName {
								if dependence.ecsID == ds.Attributes[fmt.Sprintf("volumes.%d.attachments.0.instance_id", i)] {
									return nil
								} else {
									return fmt.Errorf("绑定云主机失败")
								}
							}
						}
						return fmt.Errorf("未找到目标元素")
					},
				),
			},
			{
				ResourceName:      associationResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[associationResourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", associationResourceName)
					}
					return fmt.Sprintf("%s,%s,%s",
						rs.Primary.Attributes["ebs_id"],
						rs.Primary.Attributes["instance_id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerifyIgnore: []string{
					"master_order_id",
				},
			},
			{
				ResourceName:      associationResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[associationResourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", associationResourceName)
					}
					return fmt.Sprintf("%s,%s",
						rs.Primary.Attributes["ebs_id"],
						rs.Primary.Attributes["instance_id"],
					), nil
				},
				ImportStateVerifyIgnore: []string{
					"master_order_id",
				},
			},
			// 添加多种导入方式测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s,%s,%s",
						rs.Primary.Attributes["id"],
						rs.Primary.Attributes["project_id"],
						rs.Primary.Attributes["az_name"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"master_order_id"},
			},
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
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id", "az_name", "master_order_id"},
			},
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
						rs.Primary.Attributes["az_name"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"master_order_id"},
			},

			// 解绑
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					"sata",
					updatedSize,
					"",
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					"",
				),
			},
			// 检查解绑是否成功
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					"sata",
					updatedSize,
					"",
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					"",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					func(s *terraform.State) error {
						ds := s.RootModule().Resources[datasourceName].Primary

						count, err := strconv.Atoi(ds.Attributes["volumes.#"])
						if err != nil || count == 0 {
							return fmt.Errorf("volumes 无效: %v", ds.Attributes)
						}

						for i := 0; i < count; i++ {
							if ds.Attributes[fmt.Sprintf("volumes.%d.name", i)] == updatedName {
								if "0" == ds.Attributes[fmt.Sprintf("volumes.%d.attachments.#", i)] {
									return nil
								} else {
									return fmt.Errorf("解绑云主机失败")
								}
							}
						}
						return fmt.Errorf("未找到目标元素")
					},
				),
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					"sata",
					updatedSize,
					"",
				) + utils.LoadTestCase(
					datasourceFile, dnd,
					"",
				),
				Destroy: true,
			},
		},
	},
	)
}

func TestAccCtyunEbsMonth(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_ebs." + rnd
	resourceFile := "resource_ctyun_ebs_month.tf"
	initName := "init-ebs-" + rnd
	initSize := 60

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

			// 创建
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, initSize, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "size", strconv.Itoa(initSize)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, initName, initSize),
				Destroy: true,
			},
		},
	},
	)
}

func TestAccCtyunEbsExtra(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_ebs." + rnd
	resourceFile := "resource_ctyun_ebs.tf"
	initName := "init-ebs-extra-" + rnd
	initSize := 60

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

			// 创建
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initName,
					"xssd-0",
					initSize,
					"provisioned_iops = 1\n  delete_snap_with_ebs = true\n  multi_attach = true",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "size", strconv.Itoa(initSize)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
					resource.TestCheckResourceAttr(resourceName, "provisioned_iops", "1"),
					resource.TestCheckResourceAttr(resourceName, "delete_snap_with_ebs", "true"),
					resource.TestCheckResourceAttr(resourceName, "multi_attach", "true"),
				),
			},
			// 更新属性，同时绑定ecs
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initName,
					"xssd-0",
					initSize,
					"provisioned_iops = 2\n  delete_snap_with_ebs = false\n  multi_attach = true",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "provisioned_iops", "2"),
					resource.TestCheckResourceAttr(resourceName, "delete_snap_with_ebs", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "master_order_id"),
				),
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					initName,
					"xssd-0",
					initSize,
					"provisioned_iops = 2\n  delete_snap_with_ebs = false\n  multi_attach = true",
				),
				Destroy: true,
			},
		},
	},
	)
}
