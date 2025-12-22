package mongodb_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMongodbAccount_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_mongodb_account." + rnd
	resourceFile := "resource_ctyun_mongodb_account.tf"

	instance_id := dependence.mongodbID
	password := "@1QWs" + utils.GenerateRandomString()
	passwordUpdate := "@1qWS" + utils.GenerateRandomString()
	name := utils.GenerateRandomString()

	database := "admin"
	privileges := "readWrite"
	privilegesUpdate := "read"
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
			// 基本功能验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, instance_id, name, database, privileges, password),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", database),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, instance_id, name, database, privilegesUpdate, password),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", database),
				),
			},

			{
				Config: utils.LoadTestCase(resourceFile, rnd, instance_id, name, database, privilegesUpdate, passwordUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", database),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, instance_id, name, database, privileges, password),
				Destroy: true,
			},
		},
	})
}

func TestAccMongodbAccount_basicImportState(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_mongodb_account." + rnd
	resourceFile := "resource_ctyun_mongodb_account.tf"

	instance_id := dependence.mongodbID
	password := "@1QWs" + utils.GenerateRandomString()
	name := utils.GenerateRandomString()

	database := "admin"
	privileges := "readWrite"
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
			// 基本功能验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, instance_id, name, database, privileges, password),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", database),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					// 构造导入ID: "id,region_id"
					return fmt.Sprintf("%s,%s,%s,%s",
						rs.Primary.Attributes["name"],
						rs.Primary.Attributes["instance_id"],
						rs.Primary.Attributes["project_id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "description", "project_id", "database", "roles"},
			},
			// 3. 资源导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					// 构造导入ID: "id,region_id"
					return fmt.Sprintf("%s,%s",
						rs.Primary.Attributes["name"],
						rs.Primary.Attributes["instance_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "description", "project_id", "roles"},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, instance_id, name, database, privileges, password),
				Destroy: true,
			},
		},
	})
}

func TestMongodbAccountValidation(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceFile := "resource_ctyun_mongodb_account.tf"
	instId := dependence.mongodbID

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      utils.LoadTestCase(resourceFile, rnd, instId),
				ExpectError: regexp.MustCompile("实例名称不符合规范"),
			},
			{
				Config:      utils.LoadTestCase(resourceFile, rnd, instId),
				ExpectError: regexp.MustCompile("存在非法字符，密码仅支持大写字母、小写字母、数字和特殊字符"),
			},
		},
	})
}
