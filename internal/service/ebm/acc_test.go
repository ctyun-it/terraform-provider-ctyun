package ebm_test

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
	securityGroupID  string
	deviceType       string
	systemRaid       string
	dataRaid         string
	imageUUID        string
	ebsID            string
	ebmID            string
	securityGroupID2 string
	az2              string
}

var dependence Dependence

func TestMain(m *testing.M) {
	if skip := os.Getenv("SKIP_EBM_TEST"); skip != "" {
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
		vpcID:            outputs["vpc_id"].Value,
		subnetID:         outputs["subnet_id"].Value,
		securityGroupID:  outputs["security_group_id"].Value,
		deviceType:       outputs["device_type"].Value,
		systemRaid:       outputs["system_raid"].Value,
		dataRaid:         outputs["data_raid"].Value,
		imageUUID:        outputs["image_uuid"].Value,
		ebsID:            outputs["ebs_id"].Value,
		ebmID:            outputs["ebm_id"].Value,
		securityGroupID2: outputs["security_group_id2"].Value,
		az2:              outputs["az2"].Value,
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
