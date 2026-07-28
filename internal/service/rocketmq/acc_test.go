package rocketmq_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"os"
	"testing"
	"time"
)

const dependenceDir = "testdata/dependence"

type Dependence struct {
	regionID                 string
	vpcID                    string
	subnetID                 string
	securityGroupID          string
	instanceID               string
	exchangeName             string
	rocketmqSingleDiskType   string
	rocketmqSingleSpecName   string
	rocketmqSingleSpecName2  string
	rocketmqClusterDiskType  string
	rocketmqClusterSpecName  string
	rocketmqClusterSpecName2 string
	zone                     string
}

var dependence Dependence

func TestMain(m *testing.M) {
	os.Setenv("CTYUN_REGION_ID", "200000001703")
	// 初始化依赖资源
	fmt.Println("开始初始化依赖资源")
	outputs, err := terraform.ApplyResource(dependenceDir)
	if err != nil {
		fmt.Println(err)
		terraform.DestroyResource(dependenceDir)
		os.Exit(1)
	}

	dependence = Dependence{
		vpcID:                    outputs["vpc_id"].Value,
		subnetID:                 outputs["subnet_id"].Value,
		securityGroupID:          outputs["security_group_id"].Value,
		instanceID:               outputs["instance_id"].Value,
		exchangeName:             outputs["exchange_name"].Value,
		rocketmqClusterDiskType:  outputs["rocketmq_cluster_disk_type"].Value,
		rocketmqClusterSpecName:  outputs["rocketmq_cluster_spec_name"].Value,
		rocketmqClusterSpecName2: outputs["rocketmq_cluster_spec_name2"].Value,
		rocketmqSingleDiskType:   outputs["rocketmq_single_disk_type"].Value,
		rocketmqSingleSpecName:   outputs["rocketmq_single_spec_name"].Value,
		rocketmqSingleSpecName2:  outputs["rocketmq_single_spec_name2"].Value,
		zone:                     outputs["zone"].Value,
	}
	fmt.Println("依赖资源初始化完毕")

	// 执行测试用例
	code := m.Run()

	fmt.Println("开始清理依赖资源")
	// 清理依赖资源
	terraform.DestroyResource(dependenceDir)
	time.Sleep(3 * time.Minute)
	err = terraform.DestroyResource(dependenceDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("依赖资源清理完毕")

	os.Exit(code)
}
