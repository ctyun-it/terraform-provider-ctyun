package elb_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

func TestAccCtyunElbTarget(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_elb_target." + rnd
	resourceFile := "resource_ctyun_elb_target.tf"

	datasourceName := "data.ctyun_elb_targets." + dnd
	datasourceFile := "datasource_ctyun_elb_targets.tf"
	protocolPort := utils.GenerateRandomPort(1, 65535)
	targetGroupID := dependence.targetGroupID

	updatedProtocolPort := utils.GenerateRandomPort(1, 65535) + 1
	weight := 256
	updatedWeight := 1

	tfWeight := fmt.Sprintf(`weight=%d`, weight)
	updatedTfWeight := fmt.Sprintf(`weight=%d`, updatedWeight)

	// 代码合并后整合
	instanceType := "VM"
	instanceId := dependence.instanceID
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
			// 1. 基础功能测试
			// 1.1 Create验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, targetGroupID, instanceType, instanceId, protocolPort, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "instance_type", instanceType),
					resource.TestCheckResourceAttr(resourceName, "target_group_id", targetGroupID),
					resource.TestCheckResourceAttr(resourceName, "protocol_port", strconv.Itoa(protocolPort)),
				),
			},
			// 1.2 update验证，更新protocolPort和weight
			{
				Config: utils.LoadTestCase(resourceFile, rnd, targetGroupID, instanceType, instanceId, updatedProtocolPort, updatedTfWeight),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "instance_type", instanceType),
					resource.TestCheckResourceAttr(resourceName, "target_group_id", targetGroupID),
					resource.TestCheckResourceAttr(resourceName, "protocol_port", strconv.Itoa(updatedProtocolPort)),
					resource.TestCheckResourceAttr(resourceName, "weight", strconv.Itoa(updatedWeight)),
				),
			},
			// 1.3 datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, targetGroupID, instanceType, instanceId, updatedProtocolPort, updatedTfWeight) +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf(`ids=%s.id`, resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "elb_targets.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "elb_targets.0.instance_type", instanceType),
					resource.TestCheckResourceAttr(datasourceName, "elb_targets.0.target_group_id", targetGroupID),
					resource.TestCheckResourceAttr(datasourceName, "elb_targets.0.protocol_port", strconv.Itoa(updatedProtocolPort)),
					resource.TestCheckResourceAttr(datasourceName, "elb_targets.0.weight", strconv.Itoa(updatedWeight)),
				),
			},

			// 1.4 销毁,delete验证
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, targetGroupID, instanceType, instanceId, updatedProtocolPort, updatedTfWeight),
				Destroy: true,
			},

			// 2. 带weight创建
			// 2.1 Create
			{
				Config: utils.LoadTestCase(resourceFile, rnd, targetGroupID, instanceType, instanceId, protocolPort, tfWeight),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "instance_type", instanceType),
					resource.TestCheckResourceAttr(resourceName, "target_group_id", targetGroupID),
					resource.TestCheckResourceAttr(resourceName, "protocol_port", strconv.Itoa(protocolPort)),
					resource.TestCheckResourceAttr(resourceName, "weight", strconv.Itoa(weight)),
				),
			},
			// 2.2 update
			{
				Config: utils.LoadTestCase(resourceFile, rnd, targetGroupID, instanceType, instanceId, protocolPort, updatedTfWeight),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "instance_type", instanceType),
					resource.TestCheckResourceAttr(resourceName, "target_group_id", targetGroupID),
					resource.TestCheckResourceAttr(resourceName, "protocol_port", strconv.Itoa(protocolPort)),
					resource.TestCheckResourceAttr(resourceName, "weight", strconv.Itoa(updatedWeight)),
				),
			},
			// destroy
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, targetGroupID, instanceType, instanceId, updatedProtocolPort, updatedTfWeight),
				Destroy: true,
			},
		},
	})

}
