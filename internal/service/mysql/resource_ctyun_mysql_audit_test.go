package mysql_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
	"time"
)

func TestAccCtyunMysqlAudit_basic(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_mysql_audit." + rnd
	resourceFile := "resource_ctyun_mysql_audit.tf"

	projectID := "0"
	instID := dependence.mysqlID // MySQL实例ID

	// 测试数据
	auditSwitch := true   // 开启审计
	auditSwitch2 := false // 关闭审计

	// 等待函数
	wait30Seconds := func() {
		t.Logf("等待30秒让MySQL审计状态稳定...")
		time.Sleep(30 * time.Second)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1. 创建MySQL审计（开启审计）
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, instID, projectID,
					auditSwitch,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					// 基本属性验证
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "audit_switch", "true"),

					// 自定义验证函数
					func(s *terraform.State) error {
						_, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("resource not found: %s", resourceName)
						}
						return nil
					},
				),
				PreConfig: func() {
					wait30Seconds()
				},
			},
			// 2. 由于不支持更新，重新应用相同配置
			{
				Config: utils.LoadTestCase(
					resourceFile, rnd, instID, projectID,
					auditSwitch2,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "audit_switch", "false"),
					// 验证所有属性保持不变
					resource.TestCheckResourceAttr(resourceName, "instance_id", instID),
				),
			},
		},
	})
}
