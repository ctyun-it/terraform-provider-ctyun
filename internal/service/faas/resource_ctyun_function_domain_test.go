package faas_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunFunctionDomain_Basic(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_function_domain." + rnd
	dataSourceName := "data.ctyun_function_domains." + rnd

	resourceFile := "resource_ctyun_function_domain.tf"
	dataSourceFile := "datasource_ctyun_function_domains.tf"

	domainName := "yyy.zhangzhh13.xyz"
	description := "Terraform test function domain"

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
				// Step 1: 创建自定义域名
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					domainName,
					description,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "domain_name", domainName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "cname_check", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_status"),
				),
			},
			{
				// Step 2: 更新描述
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					domainName,
					description+" - updated",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", description+" - updated"),
				),
			},
			{
				// Step 3: 数据源查询测试
				Config: utils.LoadTestCase(dataSourceFile,
					rnd,
					domainName,
				) + "\n" + utils.LoadTestCase(resourceFile,
					rnd,
					domainName,
					description+" - updated",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "domains.#"),
				),
			},
			{
				// Step 4: 导入测试
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cname_check"},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, domainName, description),
				Destroy: true,
			},
		},
	})
}
