package redis_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunRedisAccounts(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_redis_account." + rnd
	datasourceName := "data.ctyun_redis_accounts." + dnd
	resourceFile := "resource_ctyun_redis_account.tf"
	datasourceFile := "datasource_ctyun_redis_accounts.tf"

	initName := "init_redis_account-" + rnd
	instanceId := dependence.instanceId
	initPassword := "Password_" + utils.GenerateRandomString()
	updatePassword := "Password_" + utils.GenerateRandomString()

	initPrivilege := "ro"
	updatePrivilege := "rw"

	initDescription := "Description1"
	updateDescription := "Description2"

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
				Config: utils.LoadTestCase(resourceFile, rnd, initName, instanceId, initPassword, initPrivilege, initDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
				),
			},
			// 更新
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, instanceId, updatePassword, updatePrivilege, updateDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initName),
					resource.TestCheckResourceAttr(resourceName, "password", updatePassword),
					resource.TestCheckResourceAttr(resourceName, "privilege", updatePrivilege),
					resource.TestCheckResourceAttr(resourceName, "description", updateDescription),
				),
			},
			// 查询
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, instanceId, updatePassword, updatePrivilege, updateDescription) +
					utils.LoadTestCase(datasourceFile, dnd, instanceId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "accounts.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"permission_info",
					"password",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, initName, instanceId, updatePassword, updatePrivilege, updateDescription) +
					utils.LoadTestCase(datasourceFile, dnd, instanceId),
				Destroy: true,
			},
		},
	})
}
