package scaling_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"os"
	"testing"
)

const dependenceDir = "testdata/dependence"

type Dependence struct {
	vpcID            string
	subnetID         string
	subnetID1        string
	securityGroupID  string
	loadbalancerID   string
	loadbalancerID1  string
	targetGroupID    string
	targetGroupID1   string
	securityGroupID1 string
	imageID          string
	imageID1         string
	keyPairID        string
	scalingGroupID   string
	instanceUUID     string
	instanceUUID1    string
	instanceUUID2    string
	instanceUUID3    string
	scalingConfigID  string
	scalingConfigID1 string
}

var dependence Dependence

func TestMain(m *testing.M) {
	if skip := os.Getenv("SKIP_SCALING_TEST"); skip != "" {
		return
	}
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
		subnetID1:       outputs["subnet_id1"].Value,
		securityGroupID: outputs["security_group_id"].Value,
		loadbalancerID:  outputs["elb_loadbalancer_id"].Value,
		loadbalancerID1: outputs["elb_loadbalancer_id1"].Value,
		//loadbalancerID:  "lb-qf0g35z4w5",
		//loadbalancerID1: "lb-ugjpk8itv5",
		targetGroupID:  outputs["elb_target_group_id"].Value,
		targetGroupID1: outputs["elb_target_group_id1"].Value,
		//targetGroupID:    "tg-ntp6ws9y6b",
		//targetGroupID1:   "tg-naw5sqhn8t",
		securityGroupID1: outputs["security_group_id1"].Value,
		imageID:          "e419d569-1a16-4e4e-9efc-5bc3773ca6bf",
		keyPairID:        outputs["key_pair_id"].Value,
		imageID1:         "995ecd83-c011-498b-bec7-9ab585255f9e",
		scalingGroupID:   outputs["scaling_group_id"].Value,
		instanceUUID:     outputs["instance_uuid"].Value,
		//instanceUUID: "fb081ebe-87e8-6d68-951b-d8c4a56cf4fe",
		instanceUUID1: outputs["instance_uuid1"].Value,
		//instanceUUID1:    "d8c0c1e1-3dde-9b03-b950-6bab59aa37f8",
		instanceUUID2:    outputs["instance_uuid2"].Value,
		instanceUUID3:    outputs["instance_uuid3"].Value,
		scalingConfigID:  outputs["scaling_config_id"].Value,
		scalingConfigID1: outputs["scaling_config_id1"].Value,
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
