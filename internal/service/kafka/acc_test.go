package kafka_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"os"
	"testing"
)

const dependenceDir = "testdata/dependence"

type Dependence struct {
	vpcID                 string
	subnetID              string
	securityGroupID       string
	kafkaSingleDiskType   string
	kafkaSingleSpecName   string
	kafkaSingleSpecName2  string
	kafkaClusterDiskType  string
	kafkaClusterSpecName  string
	kafkaClusterSpecName2 string
}

var dependence Dependence

func TestMain(m *testing.M) {
	if skip := os.Getenv("SKIP_KAFKA_TEST"); skip != "" {
		return
	}
	// 初始化依赖资源
	fmt.Println("开始初始化依赖资源")
	outputs, err := terraform.ApplyResource(dependenceDir)
	if err != nil {
		fmt.Println(err)
		terraform.DestroyResource(dependenceDir)
		os.Exit(1)
	}
	dependence = Dependence{
		vpcID:                 outputs["vpc_id"].Value,
		subnetID:              outputs["subnet_id"].Value,
		securityGroupID:       outputs["security_group_id"].Value,
		kafkaClusterDiskType:  outputs["kafka_cluster_disk_type"].Value,
		kafkaClusterSpecName:  outputs["kafka_cluster_spec_name"].Value,
		kafkaClusterSpecName2: outputs["kafka_cluster_spec_name2"].Value,
		kafkaSingleDiskType:   outputs["kafka_single_disk_type"].Value,
		kafkaSingleSpecName:   outputs["kafka_single_spec_name"].Value,
		kafkaSingleSpecName2:  outputs["kafka_single_spec_name2"].Value,
	}
	fmt.Println("依赖资源初始化完毕")

	// 执行测试用例
	code := m.Run()

	fmt.Println("开始清理依赖资源")
	// 清理依赖资源
	terraform.DestroyResource(dependenceDir)
	fmt.Println("依赖资源清理完毕")

	os.Exit(code)
}
