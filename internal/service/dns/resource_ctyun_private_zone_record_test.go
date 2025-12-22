package dns_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

// 测试用例1: 基础A记录测试
func TestAccCtyunPrivateZoneRecord_A(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_private_zone_record." + rnd
	resourceFile := "resource_ctyun_private_zone_record.tf"
	datasourceFile := "datasource_ctyun_private_zone_records.tf"
	datasourceName := "data.ctyun_private_zone_records." + dnd

	// 测试数据
	recordName := "test." + rnd
	zoneID := dependence.zoneID // 假设在依赖中定义了zoneID
	values := `"192.168.1.1","192.168.1.2"`
	updatedValues := `"192.168.1.3"`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建A记录
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneID,
					"A",
					values,
					300, // ttl
					recordName,
					"Test A record", // description
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "zone_id", zoneID),
					resource.TestCheckResourceAttr(resourceName, "type", "A"),
					resource.TestCheckResourceAttr(resourceName, "name", recordName),
					resource.TestCheckResourceAttr(resourceName, "ttl", "300"),
					resource.TestCheckResourceAttr(resourceName, "description", "Test A record"),
					resource.TestCheckResourceAttr(resourceName, "value_list.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "value_list.*", "192.168.1.1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "value_list.*", "192.168.1.2"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "update_time"),
				),
			},
			// 2. 更新测试 - 修改值和描述
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneID,
					"A",
					updatedValues, // 更新值
					600,           // 更新ttl
					recordName,
					"Updated A record", // 更新描述
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "value_list.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "value_list.*", "192.168.1.3"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "600"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated A record"),
				),
			},
			// datasource验证
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneID,
					"A",
					updatedValues, // 更新值
					600,           // 更新ttl
					recordName,
					"Updated A record", // 更新描述
				) + utils.LoadTestCase(
					datasourceFile, dnd, recordName, zoneID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "records.#", "1"),
				),
			},
			// import state 验证
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
				ImportStateVerifyIgnore: []string{"enabled"}, // 可选忽略
			},
			// import state 验证 - 仅ID
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
				ImportStateVerifyIgnore: []string{"enabled", "region_id"}, // 可选忽略
			},
			// 3. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneID,
					"A",
					updatedValues, // 更新值
					600,           // 更新ttl
					recordName,
					"Updated A record", // 更新描述
				),
				Destroy: true,
			},
		},
	})
}

// 测试用例2: CNAME记录测试
func TestAccCtyunPrivateZoneRecord_CNAME(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_private_zone_record." + rnd
	resourceFile := "resource_ctyun_private_zone_record.tf"
	updatedResourceFile := "resource_ctyun_private_zone_record_control.tf"

	recordName := "test-cname"
	zoneID := dependence.zoneID
	updatedValues := `"updated.example.ctyun.com"`
	values := `"init.example.ctyun.com"`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建CNAME记录
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneID,
					"CNAME",
					values,
					1000,
					recordName,
					"Test CNAME record",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "type", "CNAME"),
					resource.TestCheckResourceAttr(resourceName, "name", recordName),
					resource.TestCheckResourceAttr(resourceName, "ttl", "1000"),
					resource.TestCheckResourceAttr(resourceName, "value_list.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "value_list.*", "init.example.ctyun.com"),
				),
			},
			// 更新value
			{
				Config: utils.LoadTestCase(
					updatedResourceFile, rnd,
					zoneID,
					"CNAME",
					updatedValues,
					300,
					recordName,
					"updated CNAME record",
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "type", "CNAME"),
					resource.TestCheckResourceAttr(resourceName, "value_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "300"),
					resource.TestCheckTypeSetElemAttr(resourceName, "value_list.*", "updated.example.ctyun.com"),
				),
			},
			// 2. 清理资源
			{
				Config: utils.LoadTestCase(
					updatedResourceFile, rnd,
					zoneID,
					"CNAME",
					updatedValues,
					300,
					recordName,
					"updated CNAME record",
					false,
				),
				Destroy: true,
			},
		},
	})
}

// 测试用例3: MX记录测试
func TestAccCtyunPrivateZoneRecord_MX(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_private_zone_record." + rnd
	resourceFile := "resource_ctyun_private_zone_record.tf"

	recordName := "test-mx-record"
	zoneID := dependence.zoneID
	values := `"10 mail.example.com"`
	updatedValues := `"10 qq.mail.example.com"`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建MX记录
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneID,
					"MX",
					values,
					300,
					recordName,
					"Test MX record",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "MX"),
					resource.TestCheckResourceAttr(resourceName, "value_list.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "value_list.*", "10 mail.example.com"),
				),
			},
			// 2. 更新
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneID,
					"MX",
					updatedValues,
					300,
					recordName,
					"Test MX record",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "MX"),
					resource.TestCheckResourceAttr(resourceName, "value_list.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "value_list.*", "10 qq.mail.example.com"),
				),
			},
			// 3. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneID,
					"MX",
					updatedValues,
					300,
					recordName,
					"Test MX record",
				),
				Destroy: true,
			},
		},
	})
}

// 测试用例4: AAAA记录测试
func TestAccCtyunPrivateZoneRecord_AAAA(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_private_zone_record." + rnd
	resourceFile := "resource_ctyun_private_zone_record.tf"

	recordName := "test-aaaa-record"
	zoneID := dependence.zoneID
	values := `"2001:db8::1", "2001:db8::2"`
	updatedValues := `"2001:db8::3", "2001:db8::2", "2001:db8::4", "2001:db8::5"`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建AAAA记录
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneID,
					"AAAA",
					values,
					300,
					recordName,
					"Test AAAA record",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "AAAA"),
					resource.TestCheckResourceAttr(resourceName, "value_list.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "value_list.*", "2001:db8::1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "value_list.*", "2001:db8::2"),
				),
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneID,
					"AAAA",
					updatedValues,
					1200,
					recordName,
					"Updated AAAA record",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "AAAA"),
					resource.TestCheckResourceAttr(resourceName, "value_list.#", "4"),
					resource.TestCheckTypeSetElemAttr(resourceName, "value_list.*", "2001:db8::3"),
					resource.TestCheckTypeSetElemAttr(resourceName, "value_list.*", "2001:db8::2"),
				),
			},
			// 2. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneID,
					"AAAA",
					updatedValues,
					1200,
					recordName,
					"Updated AAAA record",
				),
				Destroy: true,
			},
		},
	})
}

// 测试用例5: TXT记录测试
func TestAccCtyunPrivateZoneRecord_TXT(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_private_zone_record." + rnd
	resourceFile := "resource_ctyun_private_zone_record.tf"

	recordName := "test-txt-record"
	zoneID := dependence.zoneID
	values := `"v spf1 include spf example com all"`
	updatedValues := `"v spf1 include spf example com all", "string1", "string2"`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建TXT记录
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneID,
					"TXT",
					values,
					300,
					recordName,
					"Test TXT record",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "TXT"),
					resource.TestCheckResourceAttr(resourceName, "value_list.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "value_list.*", "v spf1 include spf example com all"),
				),
			},
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneID,
					"TXT",
					updatedValues,
					300,
					recordName,
					"Test TXT record",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "TXT"),
					resource.TestCheckResourceAttr(resourceName, "value_list.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceName, "value_list.*", "v spf1 include spf example com all"),
					resource.TestCheckTypeSetElemAttr(resourceName, "value_list.*", "string1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "value_list.*", "string2"),
				),
			},
			// 2. 清理资源
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd,
					zoneID,
					"TXT",
					updatedValues,
					300,
					recordName,
					"Test TXT record",
				),
				Destroy: true,
			},
		},
	})
}

// 测试用例6: 测试无名称记录（根域名记录）
//func TestAccCtyunPrivateZoneRecord_SRV(t *testing.T) {
//	t.Setenv("TF_ACC", "1")
//	rnd := utils.GenerateRandomString()
//	resourceName := "ctyun_private_zone_record." + rnd
//	resourceFile := "resource_ctyun_private_zone_record.tf"
//
//	recordName := "test-srv-record"
//	zoneID := dependence.zoneID
//	values := `"3 0 2176 xmpp-server.example.com"` // 优先级 权重 端口 目标主机
//	updastedValues := `"1 1 8080 updated-xmpp-server.example.com", "3 0 2176 xmpp-server.example.com"`
//
//	resource.Test(t, resource.TestCase{
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		Steps: []resource.TestStep{
//			// 1. 创建无名称的A记录（根域名）
//			{
//				Config: utils.LoadTestCase(
//					resourceFile, rnd,
//					zoneID, "SRV",
//					values, 300, recordName, "init srv record test"),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					resource.TestCheckResourceAttr(resourceName, "type", "SRV"),
//					resource.TestCheckResourceAttr(resourceName, "value_list.#", "1"),
//					resource.TestCheckTypeSetElemAttr(resourceName, "value_list.*", "3 0 2176 xmpp-server.example.com"),
//					resource.TestCheckResourceAttr(resourceName, "name", recordName),
//				),
//			},
//			// 2. 更新
//			{
//				Config: utils.LoadTestCase(
//					resourceFile, rnd,
//					zoneID, "SRV",
//					updastedValues, 300, recordName, "updated srv record test"),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					resource.TestCheckResourceAttr(resourceName, "type", "SRV"),
//					resource.TestCheckResourceAttr(resourceName, "value_list.#", "2"),
//					resource.TestCheckTypeSetElemAttr(resourceName, "value_list.*", "3 0 2176 xmpp-server.example.com"),
//					resource.TestCheckResourceAttr(resourceName, "name", recordName),
//				),
//			},
//			// 2. 清理资源
//			{
//				Config: utils.LoadTestCase(
//					resourceFile, rnd,
//					zoneID, "SRV",
//					updastedValues, 300, recordName, "updated srv record test"),
//				Destroy: true,
//			},
//		},
//	})
//}
