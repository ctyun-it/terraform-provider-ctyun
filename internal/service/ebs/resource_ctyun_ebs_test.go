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
	initName := "init-ebs"
	initSize := 60

	updatedName := "updated-ebs"
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
					initSize,
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
					updatedSize,
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
					updatedSize,
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
					updatedSize,
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
					updatedSize,
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
				ResourceName:            associationResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			// 解绑
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedName,
					updatedSize,
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
					updatedSize,
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
					updatedSize,
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
	initName := "init-ebs"
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
				Config: utils.LoadTestCase(resourceFile, rnd, initName, initSize),
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
