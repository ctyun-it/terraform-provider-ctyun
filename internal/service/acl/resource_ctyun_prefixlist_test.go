package acl_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func mapListToStr(rules []map[string]string) string {
	rulesStr := ""
	for i, rule := range rules {
		if i > 0 {
			rulesStr += ",\n"
		}
		rulesStr += fmt.Sprintf(`
		{
			cidr = "%s"
			description = "%s"
		}`, rule["cidr"], rule["description"])
	}
	return rulesStr
}

func TestAccCtyunPrefix(t *testing.T) {

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_prefix_list." + rnd
	resourceFile := "resource_ctyun_prefix.tf"

	// 测试数据
	prefixName := "test-prefix-" + rnd
	initialDescription := "Initial prefix list description"
	updatedDescription := "Updated prefix list description"
	initialRules := []map[string]string{
		{"cidr": "192.168.1.0/24", "description": "Rule 1"},
		{"cidr": "10.0.0.0/16", "description": "Rule 2"},
	}
	//updatedRules := []map[string]string{
	//	{"cidr": "172.16.0.0/16", "description": "Updated rule 1"},
	//	{"cidr": "192.168.0.0/24", "description": "Updated rule 2"},
	//	{"cidr": "10.1.0.0/16", "description": "New rule"},
	//}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建IPv4前缀列表测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					prefixName, initialDescription,
					100, "ipv4", mapListToStr(initialRules),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", prefixName),
					resource.TestCheckResourceAttr(resourceName, "description", initialDescription),
					resource.TestCheckResourceAttr(resourceName, "limit", "100"),
					resource.TestCheckResourceAttr(resourceName, "address_type", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "prefix_list_rules.#", "2"),
					//resource.TestCheckResourceAttr(resourceName, "prefix_list_rules.0.cidr", "192.168.1.0/24"),
					//resource.TestCheckResourceAttr(resourceName, "prefix_list_rules.0.description", "Rule 1"),
					//resource.TestCheckResourceAttr(resourceName, "prefix_list_rules.1.cidr", "10.0.0.0/16"),
					//resource.TestCheckResourceAttr(resourceName, "prefix_list_rules.1.description", "Rule 2"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "update_time"),
				),
			},
			// 2. 更新前缀列表测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					prefixName+"-updated", updatedDescription,
					100, "ipv4", mapListToStr(initialRules),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", prefixName+"-updated"),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
			},
			// 3. 导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s",
						rs.Primary.Attributes["id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{}, // 不可更新的属性
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
				ImportStateVerifyIgnore: []string{}, // 不可更新的属性
			},
			// 4. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					prefixName+"-updated", updatedDescription,
					150, "ipv4", mapListToStr(initialRules),
				),
				Destroy: true,
			},
		},
	})
}

// 测试用例2: IPv6前缀列表测试
func TestAccCtyunPrefixIPv6(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_prefix_list." + rnd
	resourceFile := "resource_ctyun_prefix.tf"

	datasourceName := "data.ctyun_prefix_lists." + dnd
	datasourceFile := "datasource_ctyun_prefixlists.tf"
	// 测试数据
	prefixName := "test-prefix-ipv6-" + rnd
	description := "IPv6 prefix list"
	rules := []map[string]string{
		{"cidr": "2001:db8::/32", "description": "IPv6 rule 1"},
		{"cidr": "2001:db8:1::/48", "description": "IPv6 rule 2"},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建IPv6前缀列表测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					prefixName, description,
					100, "ipv6", mapListToStr(rules),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "address_type", "ipv6"),
					resource.TestCheckResourceAttr(resourceName, "prefix_list_rules.#", "2"),
					//resource.TestCheckResourceAttr(resourceName, "prefix_list_rules.0.cidr", "2001:db8::/32"),
					//resource.TestCheckResourceAttr(resourceName, "prefix_list_rules.0.description", "IPv6 rule 1"),
					//resource.TestCheckResourceAttr(resourceName, "prefix_list_rules.1.cidr", "2001:db8:1::/48"),
					//resource.TestCheckResourceAttr(resourceName, "prefix_list_rules.1.description", "IPv6 rule 2"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. datasource验证
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					prefixName, description,
					100, "ipv6", mapListToStr(rules),
				) + utils.LoadTestCase(datasourceFile, dnd,
					fmt.Sprintf("%s.id", resourceName), prefixName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "prefix_lists.0.name", prefixName),
					resource.TestCheckResourceAttr(datasourceName, "prefix_lists.0.limit", "100"),
					resource.TestCheckResourceAttr(datasourceName, "prefix_lists.0.address_type", "ipv6"),
				),
			},
			// 3. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					prefixName, description,
					100, "ipv6", mapListToStr(rules),
				),
			},
		},
	})
}

// 测试用例3: 边界值测试
func TestAccCtyunPrefixBoundaryValues(t *testing.T) {

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_prefix_list." + rnd
	resourceFile := "resource_ctyun_prefix.tf"

	// 测试数据
	prefixName := "prefix-boundary-" + rnd
	description := "Boundary values test"
	minRules := []map[string]string{
		{"cidr": "192.168.1.0/24", "description": "Single rule"},
	}
	maxRules := make([]map[string]string, 200)
	for i := 0; i < 200; i++ {
		maxRules[i] = map[string]string{
			"cidr":        fmt.Sprintf("10.%d.0.0/16", i),
			"description": fmt.Sprintf("Rule %d", i+1),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 最小规则数测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					prefixName, description,
					1, "ipv4", mapListToStr(minRules),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "limit", "1"),
					resource.TestCheckResourceAttr(resourceName, "prefix_list_rules.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 2. 最大规则数测试
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					prefixName+"-max", description,
					200, "ipv4", mapListToStr(maxRules),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "limit", "200"),
					resource.TestCheckResourceAttr(resourceName, "prefix_list_rules.#", "200"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// 3. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					prefixName+"-max", description,
					200, "ipv4", mapListToStr(maxRules),
				),
				Destroy: true,
			},
		},
	})
}
