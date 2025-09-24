package elb_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

// 测试经典型ELB创建、升级为保障型ELB、保障型负载均衡的信息修改和退订
// 无法测试升级保障型ELB，目前各资源池既支持经典型elb，又支持保障型elb的资源池传统型ELB均售罄

func TestAccCtyunElbLoadBalancerPg(t *testing.T) {

	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_elb_loadbalancer." + rnd
	resourceFile := "resource_ctyun_elb_loadbalancer.tf"
	dnd := utils.GenerateRandomString()
	datasourceName := "data.ctyun_elb_loadbalancers." + dnd
	datasourceFile := "datasource_ctyun_elb_loadbalancers.tf"
	name := "elb_" + utils.GenerateRandomString()
	//slaName := "elb.s1.small"
	//resourceType := "external"
	resourceType := "internal"

	updateSlaName := "elb.s2.small"
	cycleType := `cycle_type="month"`
	CycleCount := `cycle_count=1`
	//eip := `eip_id="eip-wr742lk8g3"`
	eip := ""

	update2SlaName := "elb.s3.small"

	updateName := "elb_pg_" + utils.GenerateRandomString()
	updateDescription := "terraform测试——" + utils.GenerateRandomString()

	vpcID := dependence.vpcID
	subnetID := dependence.subnetID

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

			// 创建保障型elb
			{
				Config: utils.LoadTestCase(resourceFile, rnd, subnetID, name, updateSlaName, resourceType, vpcID, "", cycleType, CycleCount, eip),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "sla_name", updateSlaName),
				),
			},
			// 保障型elb变配测试
			{
				Config: utils.LoadTestCase(resourceFile, rnd, subnetID, name, update2SlaName, resourceType, vpcID, "", cycleType, CycleCount, eip),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "sla_name", update2SlaName),
				),
			},
			// 保障型elb基本信息更新测试
			{
				Config: utils.LoadTestCase(resourceFile, rnd, subnetID, updateName, update2SlaName, resourceType, vpcID, updateDescription, cycleType, CycleCount, eip),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", updateDescription),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, subnetID, updateName, update2SlaName, resourceType, vpcID, updateDescription, cycleType, CycleCount, eip) +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf(`ids=%s.id`, resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "elbs.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "elbs.0.id"),
					resource.TestCheckResourceAttr(datasourceName, "elbs.0.name", updateName),
					resource.TestCheckResourceAttr(datasourceName, "elbs.0.description", updateDescription),
					resource.TestCheckResourceAttr(datasourceName, "elbs.0.sla_name", update2SlaName),
					resource.TestCheckResourceAttr(datasourceName, "elbs.0.resource_type", resourceType),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, subnetID, updateName, update2SlaName, resourceType, vpcID, updateDescription, cycleType, CycleCount, eip),
				Destroy: true,
			},
		},
	})

}
