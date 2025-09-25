package sfs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"os"
	"testing"
)

const dependenceDir = "testdata/dependence"

type Dependence struct {
	vpcID                 string
	vpcID1                string
	subnetID              string
	SfsUID                string
	sfsPermissionGroupID  string
	sfsPermissionGroupID1 string
}

var dependence Dependence

func TestMain(m *testing.M) {
	if skip := os.Getenv("SKIP_SFS_TEST"); skip != "" {
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
		vpcID:    outputs["vpc_id"].Value,
		vpcID1:   outputs["vpc_id1"].Value,
		subnetID: outputs["subnet_id"].Value,
		SfsUID:   outputs["sfs_uid"].Value,
		//SfsUID:                "804558e0-5ad7-5c92-9103-60823c57e524",
		sfsPermissionGroupID:  outputs["sfs_permission_group_id"].Value,
		sfsPermissionGroupID1: outputs["sfs_permission_group_id1"].Value,
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
