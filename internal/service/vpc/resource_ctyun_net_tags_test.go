package vpc_test

//
//import (
//	"fmt"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
//	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
//	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
//	"github.com/hashicorp/terraform-plugin-testing/terraform"
//	"testing"
//)
//
//func TestAccNewCtyunNetTagsResource_vpc(t *testing.T) {
//	// 支持并行执行
//	t.Parallel()
//
//	rnd := utils.GenerateRandomString()
//	dnd := utils.GenerateRandomString()
//
//	resourceName := "ctyun_net_tags." + rnd
//	datasourceName := "data.ctyun_net_tagss." + dnd
//
//	resourceFile := "resource_ctyun_net_tags.tf"
//	datasourceFile := "datasource_ctyun_net_tags.tf"
//	// 使用fmt实现标签格式化
//	tags := fmt.Sprintf(`{
//  key   = "environment"
//  value = "production"
//},
//{
//  key   = "department"
//  value = "devops"
//}`)
//	updateTags := fmt.Sprintf(`{
//  key   = "environment"
//  value = "production-update"
//},
//{
//  key   = "department-update"
//  value = "devops-update"
//}`)
//
//	resource.Test(t, resource.TestCase{
//		CheckDestroy: func(s *terraform.State) error {
//			_, exists := s.RootModule().Resources[resourceName]
//			if exists {
//				return fmt.Errorf("resource destroy failed")
//			}
//			return nil
//		},
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		Steps: []resource.TestStep{
//			// 1. VPC资源标签创建验证
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeVpc, dependence.vpcID, tags),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "resource_type", business.ResourceTypeVpc),
//					resource.TestCheckResourceAttr(resourceName, "resource_id", dependence.vpcID),
//				),
//			},
//			// 2. VPC资源标签更新验证
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeVpc, dependence.vpcID, updateTags),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "resource_type", business.ResourceTypeVpc),
//					resource.TestCheckResourceAttr(resourceName, "resource_id", dependence.vpcID),
//				),
//			},
//			// 3. datasource验证
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeVpc, dependence.vpcID, updateTags) +
//					utils.LoadTestCase(datasourceFile, dnd, business.ResourceTypeVpc, dependence.vpcID),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(datasourceName, "resource_type", business.ResourceTypeVpc),
//					resource.TestCheckResourceAttr(datasourceName, "resource_id", dependence.vpcID),
//					resource.TestCheckResourceAttrSet(datasourceName, "tags.#"),
//				),
//			},
//			// 4. 销毁验证
//			{
//				Config:  utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeVpc, dependence.vpcID, updateTags),
//				Destroy: true,
//			},
//		},
//	})
//}
//
//func TestAccNewCtyunNetTagsResource_subnet(t *testing.T) {
//	// 支持并行执行
//	t.Parallel()
//
//	rnd := utils.GenerateRandomString()
//	dnd := utils.GenerateRandomString()
//
//	resourceName := "ctyun_net_tags." + rnd
//	datasourceName := "data.ctyun_net_tagss." + dnd
//
//	resourceFile := "resource_ctyun_net_tags.tf"
//	datasourceFile := "datasource_ctyun_net_tags.tf"
//	// 使用fmt实现标签格式化
//	tags := fmt.Sprintf(`{
//  key   = "subnet-env"
//  value = "test"
//}`)
//	updateTags := fmt.Sprintf(`{
//  key   = "subnet-env"
//  value = "production"
//}`)
//
//	resource.Test(t, resource.TestCase{
//		CheckDestroy: func(s *terraform.State) error {
//			_, exists := s.RootModule().Resources[resourceName]
//			if exists {
//				return fmt.Errorf("resource destroy failed")
//			}
//			return nil
//		},
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		Steps: []resource.TestStep{
//			// 1. 子网资源标签创建验证
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeSubnet, dependence.subnetID, tags),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "resource_type", business.ResourceTypeSubnet),
//					resource.TestCheckResourceAttr(resourceName, "resource_id", dependence.subnetID),
//				),
//			},
//			// 2. 子网资源标签更新验证
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeSubnet, dependence.subnetID, updateTags),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "resource_type", business.ResourceTypeSubnet),
//					resource.TestCheckResourceAttr(resourceName, "resource_id", dependence.subnetID),
//				),
//			},
//			// 3. datasource验证
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeSubnet, dependence.subnetID, updateTags) +
//					utils.LoadTestCase(datasourceFile, dnd, business.ResourceTypeSubnet, dependence.subnetID),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(datasourceName, "resource_type", business.ResourceTypeSubnet),
//					resource.TestCheckResourceAttr(datasourceName, "resource_id", dependence.subnetID),
//					resource.TestCheckResourceAttrSet(datasourceName, "tags.#"),
//				),
//			},
//			// 4. 销毁验证
//			{
//				Config:  utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeSubnet, dependence.subnetID, updateTags),
//				Destroy: true,
//			},
//		},
//	})
//}
//
//func TestAccNewCtyunNetTagsResource_securityGroup(t *testing.T) {
//	// 支持并行执行
//	t.Parallel()
//
//	rnd := utils.GenerateRandomString()
//	dnd := utils.GenerateRandomString()
//
//	resourceName := "ctyun_net_tags." + rnd
//	datasourceName := "data.ctyun_net_tagss." + dnd
//
//	resourceFile := "resource_ctyun_net_tags.tf"
//	datasourceFile := "datasource_ctyun_net_tags.tf"
//	// 使用fmt实现标签格式化
//	tags := fmt.Sprintf(`{
//  key   = "sg-type"
//  value = "test"
//}`)
//	updateTags := fmt.Sprintf(`{
//  key   = "sg-type"
//  value = "production"
//}`)
//
//	resource.Test(t, resource.TestCase{
//		CheckDestroy: func(s *terraform.State) error {
//			_, exists := s.RootModule().Resources[resourceName]
//			if exists {
//				return fmt.Errorf("resource destroy failed")
//			}
//			return nil
//		},
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		Steps: []resource.TestStep{
//			// 1. 安全组资源标签创建验证
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeSecurityGroup, dependence.securityGroupID, tags),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "resource_type", business.ResourceTypeSecurityGroup),
//					resource.TestCheckResourceAttr(resourceName, "resource_id", dependence.securityGroupID),
//				),
//			},
//			// 2. 安全组资源标签更新验证
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeSecurityGroup, dependence.securityGroupID, updateTags),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "resource_type", business.ResourceTypeSecurityGroup),
//					resource.TestCheckResourceAttr(resourceName, "resource_id", dependence.securityGroupID),
//				),
//			},
//			// 3. datasource验证
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeSecurityGroup, dependence.securityGroupID, updateTags) +
//					utils.LoadTestCase(datasourceFile, dnd, business.ResourceTypeSecurityGroup, dependence.securityGroupID),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(datasourceName, "resource_type", business.ResourceTypeSecurityGroup),
//					resource.TestCheckResourceAttr(datasourceName, "resource_id", dependence.securityGroupID),
//					resource.TestCheckResourceAttrSet(datasourceName, "tags.#"),
//				),
//			},
//			// 4. 销毁验证
//			{
//				Config:  utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeSecurityGroup, dependence.securityGroupID, updateTags),
//				Destroy: true,
//			},
//		},
//	})
//}
//
//func TestAccNewCtyunNetTagsResource_eip(t *testing.T) {
//	// 支持并行执行
//	t.Parallel()
//
//	rnd := utils.GenerateRandomString()
//	dnd := utils.GenerateRandomString()
//
//	resourceName := "ctyun_net_tags." + rnd
//	datasourceName := "data.ctyun_net_tagss." + dnd
//
//	resourceFile := "resource_ctyun_net_tags.tf"
//	datasourceFile := "datasource_ctyun_net_tags.tf"
//	// 使用fmt实现标签格式化
//	tags := fmt.Sprintf(`{
//  key   = "eip-type"
//  value = "test"
//}`)
//	updateTags := fmt.Sprintf(`{
//  key   = "eip-type"
//  value = "production"
//}`)
//
//	resource.Test(t, resource.TestCase{
//		CheckDestroy: func(s *terraform.State) error {
//			_, exists := s.RootModule().Resources[resourceName]
//			if exists {
//				return fmt.Errorf("resource destroy failed")
//			}
//			return nil
//		},
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		Steps: []resource.TestStep{
//			// 1. 弹性IP资源标签创建验证
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeEip, dependence.eipID, tags),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "resource_type", business.ResourceTypeEip),
//					resource.TestCheckResourceAttr(resourceName, "resource_id", dependence.eipID),
//				),
//			},
//			// 2. 弹性IP资源标签更新验证
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeEip, dependence.eipID, updateTags),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "resource_type", business.ResourceTypeEip),
//					resource.TestCheckResourceAttr(resourceName, "resource_id", dependence.eipID),
//				),
//			},
//			// 3. datasource验证
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeEip, dependence.eipID, updateTags) +
//					utils.LoadTestCase(datasourceFile, dnd, business.ResourceTypeEip, dependence.eipID),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(datasourceName, "resource_type", business.ResourceTypeEip),
//					resource.TestCheckResourceAttr(datasourceName, "resource_id", dependence.eipID),
//					resource.TestCheckResourceAttrSet(datasourceName, "tags.#"),
//				),
//			},
//			// 4. 销毁验证
//			{
//				Config:  utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeEip, dependence.eipID, updateTags),
//				Destroy: true,
//			},
//		},
//	})
//}
//
//func TestAccNewCtyunNetTagsResource_bandwidth(t *testing.T) {
//	// 支持并行执行
//	t.Parallel()
//
//	rnd := utils.GenerateRandomString()
//	dnd := utils.GenerateRandomString()
//
//	resourceName := "ctyun_net_tags." + rnd
//	datasourceName := "data.ctyun_net_tagss." + dnd
//
//	resourceFile := "resource_ctyun_net_tags.tf"
//	datasourceFile := "datasource_ctyun_net_tags.tf"
//	// 使用fmt实现标签格式化
//	tags := fmt.Sprintf(`{
//  key   = "bandwidth-type"
//  value = "test"
//}`)
//	updateTags := fmt.Sprintf(`{
//  key   = "bandwidth-type"
//  value = "production"
//}`)
//
//	resource.Test(t, resource.TestCase{
//		CheckDestroy: func(s *terraform.State) error {
//			_, exists := s.RootModule().Resources[resourceName]
//			if exists {
//				return fmt.Errorf("resource destroy failed")
//			}
//			return nil
//		},
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		Steps: []resource.TestStep{
//			// 1. 带宽资源标签创建验证
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeBandwidth, dependence.bandwidthID, tags),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "resource_type", business.ResourceTypeBandwidth),
//					resource.TestCheckResourceAttr(resourceName, "resource_id", dependence.bandwidthID),
//				),
//			},
//			// 2. 带宽资源标签更新验证
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeBandwidth, dependence.bandwidthID, updateTags),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "resource_type", business.ResourceTypeBandwidth),
//					resource.TestCheckResourceAttr(resourceName, "resource_id", dependence.bandwidthID),
//				),
//			},
//			// 3. datasource验证
//			{
//				Config: utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeBandwidth, dependence.bandwidthID, updateTags) +
//					utils.LoadTestCase(datasourceFile, dnd, business.ResourceTypeBandwidth, dependence.bandwidthID),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr(datasourceName, "resource_type", business.ResourceTypeBandwidth),
//					resource.TestCheckResourceAttr(datasourceName, "resource_id", dependence.bandwidthID),
//					resource.TestCheckResourceAttrSet(datasourceName, "tags.#"),
//				),
//			},
//			// 4. 销毁验证
//			{
//				Config:  utils.LoadTestCase(resourceFile, rnd, business.ResourceTypeBandwidth, dependence.bandwidthID, updateTags),
//				Destroy: true,
//			},
//		},
//	})
//}
