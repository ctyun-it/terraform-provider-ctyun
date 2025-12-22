package dns_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"os"
	"testing"
)

const dependenceDir = "testdata/dependence"

type Dependence struct {
	vpcID  string
	vpcID1 string
	vpcID2 string
	vpcID3 string
	vpcID4 string
	vpcID5 string
	zoneID string
}

var dependence Dependence

func TestMain(m *testing.M) {
	if skip := os.Getenv("SKIP_PGSQL_TEST"); skip != "" {
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
		vpcID:  outputs["vpc_id"].Value,
		vpcID1: outputs["vpc_id1"].Value,
		vpcID2: outputs["vpc_id2"].Value,
		vpcID3: outputs["vpc_id3"].Value,
		vpcID4: outputs["vpc_id4"].Value,
		vpcID5: outputs["vpc_id5"].Value,
		zoneID: outputs["zone_id"].Value,
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
