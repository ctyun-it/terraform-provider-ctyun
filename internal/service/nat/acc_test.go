package nat_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"os"
	"testing"
)

const dependenceDir = "testdata/dependence"

type Dependence struct {
	vpcID     string
	eipID     string
	eipID1    string
	natID     string
	subnetID1 string
	subnetID2 string
	ecsID     string
}

var dependence Dependence

func TestMain(m *testing.M) {
	fmt.Println("开始初始化依赖资源")
	outputs, err := terraform.ApplyResource(dependenceDir)
	if err != nil {
		fmt.Println(err)
		terraform.DestroyResource(dependenceDir)
		os.Exit(1)
	}
	dependence = Dependence{
		vpcID:     outputs["vpc_id"].Value,
		eipID:     outputs["eip_id"].Value,
		eipID1:    outputs["eip_id1"].Value,
		natID:     outputs["nat_id"].Value,
		subnetID1: outputs["subnet_id1"].Value,
		subnetID2: outputs["subnet_id2"].Value,
		ecsID:     outputs["ecs_id"].Value,
	}

	fmt.Println("依赖资源初始化完毕")

	// 执行测试用例
	code := m.Run()
	fmt.Println("开始清理依赖资源")
	// 清理依赖资源
	err = terraform.DestroyResource(dependenceDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("依赖资源清理完毕")
	os.Exit(code)
}
