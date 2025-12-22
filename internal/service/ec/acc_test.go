package ec_test

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
	subnetID2        string
	expressConnectID string
	cloudGatewayId   string
	vpcInstanceVpcID string
	rtbID            string
	regionPeerID     string
	cgwID1           string
	cgwID2           string
	packetID         string
	sdwanId          string
}

var dependence Dependence

func TestMain(m *testing.M) {
	// 初始化依赖资源
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
		expressConnectID: outputs["ctyun_express_connect_id"].Value,
		cloudGatewayId:   outputs["ctyun_ec_cloud_gateway_id"].Value,
		vpcID:            outputs["vpc_id"].Value,
		subnetID:         outputs["subnet_id"].Value,
		subnetID2:        outputs["subnet_id2"].Value,
		vpcInstanceVpcID: outputs["vpc_instance_vpc_id"].Value,
		rtbID:            outputs["rtb_id"].Value,
		regionPeerID:     outputs["region_peer_id"].Value,
		//regionPeerID: "",
		cgwID1:   outputs["cgw_id1"].Value,
		cgwID2:   outputs["cgw_id2"].Value,
		packetID: outputs["packet_id"].Value,
		sdwanId:  outputs["sdwan_id"].Value,
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
