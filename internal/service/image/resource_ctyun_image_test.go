package image_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunImage_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_image." + rnd
	resourceFile := "resource_ctyun_image.tf"

	imageName := "tf-image-test-" + utils.GenerateRandomString()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				// 测试创建
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					imageName,
					"Ubuntu",
					"22.04",
					"https://jiangsu-10.zos.ctyun.cn/bucket-29bf/test-image", // file_source - 实际测试时需要替换为有效值
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", imageName),
					resource.TestCheckResourceAttr(resourceName, "os_distro", "Ubuntu"),
					resource.TestCheckResourceAttr(resourceName, "os_version", "22.04"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				// 测试导入
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("not found: %s", resourceName)
					}
					return fmt.Sprintf("%s", rs.Primary.ID), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"file_source", // file_source 在导入时不设置
				},
			},
			{
				// 测试导入 - 包含 region_id
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("not found: %s", resourceName)
					}
					regionId := rs.Primary.Attributes["region_id"]
					if regionId == "" {
						return "", fmt.Errorf("region_id is not set")
					}
					return fmt.Sprintf("%s,%s", rs.Primary.ID, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"file_source", // file_source 在导入时不设置
				},
			},
			{
				// 测试销毁
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					imageName,
					"Ubuntu",
					"22.04",
					"https://jiangsu-10.zos.ctyun.cn/bucket-29bf/test-image", // file_source - 实际测试时需要替换为有效值
				),
				Destroy: true,
			},
		},
	})
}
