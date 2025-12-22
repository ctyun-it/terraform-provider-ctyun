package redis_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"os"
	"testing"
)

const dependenceDir = "testdata/dependence"

type Dependence struct {
	vpcID              string
	subnetID           string
	securityGroupID    string
	eipAddress         string
	eipID              string
	redisVersion       string
	redisEngineEdition string
	instanceId         string
	instance2Id        string
	address            string
	address2           string
	userName           string
	userPassword       string
	user2Name          string
	user2Password      string
}

var dependence Dependence

func TestMain(m *testing.M) {
	// 初始化依赖资源
	if skip := os.Getenv("SKIP_REDIS_TEST"); skip != "" {
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
		vpcID:              outputs["vpc_id"].Value,
		subnetID:           outputs["subnet_id"].Value,
		securityGroupID:    outputs["security_group_id"].Value,
		eipAddress:         outputs["eip_address"].Value,
		eipID:              outputs["eip_id"].Value,
		redisVersion:       outputs["redis_version"].Value,
		redisEngineEdition: outputs["redis_engine_edition"].Value,
		instanceId:         outputs["redis_instance_id"].Value,
		instance2Id:        outputs["redis_instance2_id"].Value,
		address:            outputs["redis_address"].Value,
		address2:           outputs["redis2_address"].Value,
		userName:           outputs["instance_account_name"].Value,
		userPassword:       outputs["instance_account_pswd"].Value,
		user2Name:          outputs["instance2_account_name"].Value,
		user2Password:      outputs["instance2_account_pswd"].Value,
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
