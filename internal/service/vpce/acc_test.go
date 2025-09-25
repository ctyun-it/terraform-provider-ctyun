package vpce_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"os"
	"testing"
)

const dependenceDir = "testdata/dependence"

type Dependence struct {
	ecsID                string
	ecsID2               string
	vpcID                string
	subnetID             string
	vpceServiceID        string
	reverseVpceServiceID string
	vpceID               string
	transitIP            string
	targetIP             string
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
		ecsID:                outputs["ecs_id"].Value,
		ecsID2:               outputs["ecs_id2"].Value,
		vpcID:                outputs["vpc_id"].Value,
		subnetID:             outputs["subnet_id"].Value,
		vpceServiceID:        outputs["vpce_service_id"].Value,
		reverseVpceServiceID: outputs["reverse_vpce_service_id"].Value,
		vpceID:               outputs["vpce_id"].Value,
		transitIP:            outputs["vpce_service_transit_ip"].Value,
		targetIP:             outputs["ecs_fixed_ip"].Value,
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
