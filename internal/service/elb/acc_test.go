package elb_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"os"
	"testing"
)

const dependenceDir = "testdata/dependence"

type Dependence struct {
	vpcID          string
	subnetID       string
	loadBalanceID  string
	loadBalanceID2 string
	healthCheckID  string
	targetGroupID  string
	targetGroupID2 string
	targetGroupID3 string
	targetGroupID4 string
	listenerID     string
	instanceID     string
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
		vpcID:    outputs["vpc_id"].Value,
		subnetID: outputs["subnet_id"].Value,
		//loadBalanceID: outputs["loadbalancer_id"].Value,
		loadBalanceID:  "",
		loadBalanceID2: outputs["loadbalancer_id_rule"].Value,
		//loadBalanceID2: "",
		healthCheckID:  outputs["health_check_id"].Value,
		targetGroupID:  outputs["target_group_id"].Value,
		targetGroupID2: outputs["target_group_id2"].Value,
		targetGroupID3: outputs["target_group_id3"].Value,
		targetGroupID4: outputs["target_group_id4"].Value,
		listenerID:     outputs["listener_id"].Value,
		//listenerID: "",
		instanceID: outputs["instance_id"].Value,
		//instanceID: "",
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
