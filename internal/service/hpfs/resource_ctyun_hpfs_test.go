package hpfs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
	"time"
)

func TestAccCtyunHpfs(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_hpfs." + rnd
	resourceFile := "resource_ctyun_hpfs.tf"

	datasourceName := "data.ctyun_hpfs_instances." + dnd
	datasourceFile := "datasource_ctyun_hfps_instances.tf"
	//
	//vpcID := dependence.vpcID
	//subnetID := dependence.subnetID
	sfsProtocol := "hpfs"
	cycleType := "on_demand"
	sfsName := "hpfs-" + utils.GenerateRandomString()
	updatedSfsName := "hpfs-" + utils.GenerateRandomString() + "-new"
	sfsSize := 512
	updatedSfsSize := 1024
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
			// 开通hpfs，
			{
				Config: utils.LoadTestCase(resourceFile, rnd, sfsProtocol, cycleType, sfsName, sfsSize),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "sfs_protocol", sfsProtocol),
					resource.TestCheckResourceAttr(resourceName, "name", sfsName),
					resource.TestCheckResourceAttr(resourceName, "sfs_size", strconv.Itoa(sfsSize)),
				),
			},
			// 变配sfs name 和 SIZE规格 512->1024
			{
				Config: utils.LoadTestCase(resourceFile, rnd, sfsProtocol, cycleType, updatedSfsName, updatedSfsSize),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "sfs_protocol", sfsProtocol),
					resource.TestCheckResourceAttr(resourceName, "name", updatedSfsName),
					resource.TestCheckResourceAttr(resourceName, "sfs_size", strconv.Itoa(updatedSfsSize)),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, sfsProtocol, cycleType, updatedSfsName, updatedSfsSize) +
					utils.LoadTestCase(datasourceFile, dnd, "available", sfsProtocol, "cn-huadong1-jsnj1A-public-ctcloud"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "hpfs_instances.0.sfs_protocol", sfsProtocol),
					resource.TestCheckResourceAttr(datasourceName, "hpfs_instances.0.sfs_status", "available"),
					resource.TestCheckResourceAttr(datasourceName, "hpfs_instances.0.az_name", "cn-huadong1-jsnj1A-public-ctcloud"),
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(30 * time.Second)
							return nil
						},
					),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, sfsProtocol, cycleType, updatedSfsName, updatedSfsSize),
				Destroy: true,
			},
		},
	})
}

// 指定AZ，指定集群和baseline
func TestAccCtyunHpfs1(t *testing.T) {
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_hpfs." + rnd
	resourceFile := "resource_ctyun_hpfs1.tf"
	//
	//vpcID := dependence.vpcID
	//subnetID := dependence.subnetID
	sfsProtocol := "hpfs"
	cluster := "hdRoce01"
	baseline := "200"
	sfsName := "hpfs-" + utils.GenerateRandomString()
	updatedSfsName := "hpfs-" + utils.GenerateRandomString() + "-new"
	azName := "cn-huadong1-jsnj1A-public-ctcloud"
	sfsSize := 512
	updatedSfsSize := 512
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
			// 开通hpfs，
			{
				Config: utils.LoadTestCase(resourceFile, rnd, sfsProtocol, sfsName, sfsSize, azName, cluster, baseline),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "sfs_protocol", sfsProtocol),
					resource.TestCheckResourceAttr(resourceName, "name", sfsName),
					resource.TestCheckResourceAttr(resourceName, "sfs_size", strconv.Itoa(sfsSize)),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", cluster),
					resource.TestCheckResourceAttr(resourceName, "baseline", baseline),
				),
			},
			// 变配sfs name 和 SIZE规格 512->1024
			{
				Config: utils.LoadTestCase(resourceFile, rnd, sfsProtocol, updatedSfsName, updatedSfsSize, azName, cluster, baseline),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "sfs_protocol", sfsProtocol),
					resource.TestCheckResourceAttr(resourceName, "name", updatedSfsName),
					resource.TestCheckResourceAttr(resourceName, "sfs_size", strconv.Itoa(updatedSfsSize)),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", cluster),
					resource.TestCheckResourceAttr(resourceName, "baseline", baseline),
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(30 * time.Second)
							return nil
						},
					),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, sfsProtocol, sfsName, sfsSize, azName, cluster, baseline),
				Destroy: true,
			},
		},
	})
}

type CtyunHpfsInstancesModel struct {
	SfsName       string   `tfsdk:"name"`            // 并行文件命名
	SfsID         string   `tfsdk:"sfs_id"`          // 并行文件唯一ID
	SfsSize       int32    `tfsdk:"sfs_size"`        // 大小(GB)
	SfsType       string   `tfsdk:"sfs_type"`        // 文件系统类型
	SfsProtocol   string   `tfsdk:"sfs_protocol"`    // 挂载协议
	SfsStatus     string   `tfsdk:"sfs_status"`      // 文件系统状态
	UsedSize      int32    `tfsdk:"used_size"`       // 已用大小(MB)
	CreateTime    int64    `tfsdk:"create_time"`     // 创建时间戳(毫秒)
	UpdateTime    int64    `tfsdk:"update_time"`     // 更新时间戳(毫秒)
	ProjectID     string   `tfsdk:"project_id"`      // 企业项目ID
	OnDemand      bool     `tfsdk:"on_demand"`       // 是否按需订购
	RegionID      string   `tfsdk:"region_id"`       // 资源池ID
	AzName        string   `tfsdk:"az_name"`         // 可用区名称
	ClusterName   string   `tfsdk:"cluster_name"`    // 集群名称
	Baseline      string   `tfsdk:"baseline"`        // 性能基线(MB/s/TB)
	HpfsSharePath string   `tfsdk:"hpfs_share_path"` // HPFS共享路径
	SecretKey     string   `tfsdk:"secret_key"`      // HPC挂载密钥
	DataflowList  []string `tfsdk:"dataflow_list"`   // 数据流动策略ID列表
	DataflowCount int32    `tfsdk:"dataflow_count"`  // 数据流动策略数量
}
