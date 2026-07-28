package cloud_assistant_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
)

const dependenceDir = "testdata/dependence"

type Dependence struct {
	ecsID1    string
	ecsID2    string
	commandID string
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
		ecsID1:    outputs["ecs_id1"].Value,
		ecsID2:    outputs["ecs_id2"].Value,
		commandID: outputs["command_id"].Value,
	}
	fmt.Println("依赖资源初始化完毕")

	code := m.Run()

	fmt.Println("开始清理依赖资源")
	err = terraform.DestroyResource(dependenceDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("依赖资源清理完毕")

	os.Exit(code)
}
