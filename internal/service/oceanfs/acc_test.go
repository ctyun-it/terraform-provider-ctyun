package oceanfs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"os"
	"testing"
)

const dependenceDir = "testdata/dependence"

type Dependence struct {
	vpcID              string
	vpcID1             string
	subnetID           string
	subnetID1          string
	permissionGroupID  string
	permissionGroupID1 string
	oceanfsID          string
}

var dependence Dependence

func TestMain(m *testing.M) {
	if skip := os.Getenv("SKIP_OCEANFS_TEST"); skip != "" {
		return
	}
	os.Setenv("CTYUN_REGION_ID", "200000003573")
	os.Setenv("CTYUN_AZ_NAME", "cn-nm-het3-1a-public-ctcloud")
	fmt.Println("开始初始化依赖资源")
	outputs, err := terraform.ApplyResource(dependenceDir)
	if err != nil {
		fmt.Println(err)
		terraform.DestroyResource(dependenceDir)
		os.Exit(1)
	}
	dependence = Dependence{
		vpcID:              outputs["vpc_id"].Value,
		vpcID1:             outputs["vpc_id1"].Value,
		subnetID:           outputs["subnet_id"].Value,
		subnetID1:          outputs["subnet_id1"].Value,
		permissionGroupID:  outputs["permission_group_id"].Value,
		permissionGroupID1: outputs["permission_group_id1"].Value,
		oceanfsID:          outputs["oceanfs_id"].Value,
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
