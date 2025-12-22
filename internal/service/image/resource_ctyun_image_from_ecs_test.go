package image_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunImageFromEcs_basic(t *testing.T) {
	//t.Parallel()

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_image_from_ecs." + rnd
	resourceFile := "ctyun_image_from_ecs_system_disk.tf"

	imageName := "tf-image-test-" + utils.GenerateRandomString()
	updatedImageName := "tf-image-test-updated-" + utils.GenerateRandomString()

	description := "测试镜像描述"
	updatedDescription := description + "-updated"

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
					description,
					dependence.instanceID,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_name", imageName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				// 测试更新
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedImageName,
					updatedDescription,
					dependence.instanceID,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_name", updatedImageName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
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
					"instance_id",
					"data_disk_id",
					"repository_id",
					"snapshot_id",
					"image_type",
					"labels",
				},
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
					regionId := rs.Primary.Attributes["region_id"]
					if regionId == "" {
						return "", fmt.Errorf("region_id is not set")
					}
					return fmt.Sprintf("%s,%s", rs.Primary.ID, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_id",
					"data_disk_id",
					"repository_id",
					"snapshot_id",
					"image_type",
					"labels",
				},
			},
			{
				// 测试销毁
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedImageName,
					updatedDescription,
					dependence.instanceID,
				),
				Destroy: true,
			},
		},
	})
}
func TestAccCtyunImageFromEcsSystemDisk_case1(t *testing.T) {
	//t.Parallel()

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_image_from_ecs." + rnd
	resourceFile := "ctyun_image_from_ecs_system_disk_case_1.tf"

	imageName := "tf-image-test-case1" + utils.GenerateRandomString()
	updatedImageName := "tf-image-updated-" + utils.GenerateRandomString()

	description := "测试镜像描述-case1"
	updatedDescription := description + "-updated"

	// 定义内存参数
	minimumRam := "2"
	minimumRam_update := "1"
	maximumRam := "16"
	maximumRam_update := "32"

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
					description,
					dependence.instanceID,
					minimumRam, // minimum_ram
					maximumRam, // maximum_ram
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_name", imageName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "minimum_ram", minimumRam),
					resource.TestCheckResourceAttr(resourceName, "maximum_ram", maximumRam),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				// 测试更新
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedImageName,
					updatedDescription,
					dependence.instanceID,
					minimumRam_update, // minimum_ram
					maximumRam_update, // maximum_ram
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_name", updatedImageName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttr(resourceName, "minimum_ram", minimumRam_update),
					resource.TestCheckResourceAttr(resourceName, "maximum_ram", maximumRam_update),
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

					regionId := rs.Primary.Attributes["region_id"]
					if regionId == "" {
						return "", fmt.Errorf("region_id is not set")
					}

					return fmt.Sprintf("%s,%s", rs.Primary.ID, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_id",
					"data_disk_id",
					"repository_id",
					"snapshot_id",
					"image_type",
					"labels",
				},
			},
			{
				// 测试销毁
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedImageName,
					updatedDescription,
					dependence.instanceID,
					minimumRam, // minimum_ram
					maximumRam, // maximum_ram
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunImageFromEcs_dataDisk(t *testing.T) {
	//t.Parallel()

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_image_from_ecs." + rnd
	resourceFile := "ctyun_image_from_ecs_data_disk.tf"

	imageName := "tf-image-data-" + utils.GenerateRandomString()
	updatedImageName := "tf-image-data-updated-" + utils.GenerateRandomString()

	description := "测试数据盘创建镜像描述"
	updatedDescription := description + "-updated"

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
				// 测试数据盘创建镜像
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					imageName,
					description,
					dependence.instanceID,
					dependence.dataDiskID,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_name", imageName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				// 测试更新
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedImageName,
					updatedDescription,
					dependence.instanceID,
					dependence.dataDiskID,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_name", updatedImageName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
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

					regionId := rs.Primary.Attributes["region_id"]
					if regionId == "" {
						return "", fmt.Errorf("region_id is not set")
					}

					return fmt.Sprintf("%s,%s", rs.Primary.ID, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_id",
					"data_disk_id",
					"repository_id",
					"snapshot_id",
					"image_type",
					"labels",
				},
			},
			{
				// 测试销毁
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedImageName,
					updatedDescription,
					dependence.instanceID,
					dependence.dataDiskID,
				),
				Destroy: true,
			},
		},
	})
}
func TestAccCtyunImageFromEcs_dataDisk_case1(t *testing.T) {
	//t.Parallel()

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_image_from_ecs." + rnd
	resourceFile := "ctyun_image_from_ecs_data_disk_case1.tf"

	imageName := "im-data-1-" + utils.GenerateRandomString()
	updatedImageName := "im-data-1-updated-" + utils.GenerateRandomString()

	description := "测试数据盘创建镜像描述"
	updatedDescription := description + "-updated"

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
				// 测试数据盘创建镜像
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					imageName,
					description,
					dependence.instanceID,
					dependence.dataDiskID,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_name", imageName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				// 测试更新
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedImageName,
					updatedDescription,
					dependence.instanceID,
					dependence.dataDiskID,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_name", updatedImageName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
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

					regionId := rs.Primary.Attributes["region_id"]
					if regionId == "" {
						return "", fmt.Errorf("region_id is not set")
					}

					return fmt.Sprintf("%s,%s", rs.Primary.ID, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_id",
					"data_disk_id",
					"repository_id",
					"snapshot_id",
					"image_type",
					"labels",
				},
			},
			{
				// 测试销毁
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedImageName,
					updatedDescription,
					dependence.instanceID,
					dependence.dataDiskID,
				),
				Destroy: true,
			},
		},
	})
}
func TestAccCtyunImageFromEcs_entire(t *testing.T) {
	//t.Parallel()

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_image_from_ecs." + rnd
	resourceFile := "ctyun_image_from_ecs_entire.tf"

	imageName := "tf-image-test-" + utils.GenerateRandomString()
	updatedImageName := "tf-image-test-updated-" + utils.GenerateRandomString()

	description := "测试整机镜像描述"
	updatedDescription := description + "-updated"

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
					description,
					dependence.instanceID,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_name", imageName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				// 测试更新
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedImageName,
					updatedDescription,
					dependence.instanceID,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_name", updatedImageName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
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

					regionId := rs.Primary.Attributes["region_id"]
					if regionId == "" {
						return "", fmt.Errorf("region_id is not set")
					}

					return fmt.Sprintf("%s,%s", rs.Primary.ID, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_id",
					"data_disk_id",
					"repository_id",
					"snapshot_id",
					"image_type",
					"labels",
				},
			},
			{
				// 测试销毁
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedImageName,
					updatedDescription,
					dependence.instanceID,
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunImageFromEcs_entire_case1(t *testing.T) {
	//t.Parallel()

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_image_from_ecs." + rnd
	resourceFile := "ctyun_image_from_ecs_entire_case1.tf"

	imageName := "im-entire-1-" + utils.GenerateRandomString()
	updatedImageName := "im-entire-updated-" + utils.GenerateRandomString()

	description := "测试整机镜像描述"
	updatedDescription := description + "-updated"

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
					description,
					dependence.instanceID,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_name", imageName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				// 测试更新
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedImageName,
					updatedDescription,
					dependence.instanceID,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_name", updatedImageName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
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

					regionId := rs.Primary.Attributes["region_id"]
					if regionId == "" {
						return "", fmt.Errorf("region_id is not set")
					}

					return fmt.Sprintf("%s,%s", rs.Primary.ID, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_id",
					"data_disk_id",
					"repository_id",
					"snapshot_id",
					"image_type",
					"labels",
				},
			},
			{
				// 测试销毁
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					updatedImageName,
					updatedDescription,
					dependence.instanceID,
				),
				Destroy: true,
			},
		},
	})
}
