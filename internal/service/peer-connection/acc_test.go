package peer_connection_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"os"
	"testing"
)

const dependenceDir = "testdata/dependence"

type Dependence struct {
	vpcID1 string
	vpcID2 string
	//crossAccountVpcID string
	//crossAccountEmail string
	peerConnectionID string
	rtbID            string
}

var dependence Dependence

func TestMain(m *testing.M) {
	if skip := os.Getenv("SKIP_PEER_CONNECT_TEST"); skip != "" {
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
		vpcID1: outputs["vpc_id"].Value,
		vpcID2: outputs["vpc_id1"].Value,
		//crossAccountVpcID: outputs["vpc_id2"].Value,
		//crossAccountEmail: "925415014@qq.com",
		peerConnectionID: outputs["peer_connection_id"].Value,
		rtbID:            outputs["rtb_id"].Value,
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
