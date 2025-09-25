package ecs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"os"
	"testing"
)

const dependenceDir = "testdata/dependence"

type Dependence struct {
	vpcID           string
	subnetID        string
	imageID         string
	flavorID        string
	flavorID2       string
	affinityGroupID string
	keyPairName     string
	keyPairName2    string
	securityGroupID string
	ecsID           string
}

var dependence Dependence

func TestMain(m *testing.M) {
	// 初始化依赖资源
	fmt.Println("开始初始化依赖资源")
	outputs, err := terraform.ApplyResource(dependenceDir)
	if err != nil {
		fmt.Println(err)
		terraform.DestroyResource(dependenceDir)
		os.Exit(1)
	}
	dependence = Dependence{
		vpcID:           outputs["vpc_id"].Value,
		subnetID:        outputs["subnet_id"].Value,
		imageID:         outputs["image_id"].Value,
		flavorID:        outputs["flavor_id"].Value,
		flavorID2:       outputs["flavor_id2"].Value,
		affinityGroupID: outputs["affinity_group_id"].Value,
		keyPairName:     outputs["key_pair_name"].Value,
		keyPairName2:    outputs["key_pair_name2"].Value,
		securityGroupID: outputs["security_group_id"].Value,
		ecsID:           outputs["ecs_id"].Value,
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
