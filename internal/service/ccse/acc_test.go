package ccse_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"os"
	"testing"
)

const dependenceDir = "testdata/dependence"

type Dependence struct {
	vpcID           string
	subnetID        string
	flavorName      string
	clusterID       string
	chartName       string
	chartVersion1   string
	chartVersion2   string
	chartValuesYaml string
	chartValuesJson string
	ecsID           string
	ebmID           string
	ecsMirrorID     string
	ebmMirrorID     string
	ebmMirrorName   string
	deviceType      string
	ebmAz           string
}

var dependence Dependence

func TestMain(m *testing.M) {
	if skip := os.Getenv("SKIP_CCSE_TEST"); skip != "" {
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
		vpcID:           outputs["vpc_id"].Value,
		subnetID:        outputs["subnet_id"].Value,
		flavorName:      outputs["flavor_name"].Value,
		clusterID:       outputs["cluster_id"].Value,
		chartName:       outputs["chart_name"].Value,
		chartVersion1:   outputs["chart_version1"].Value,
		chartVersion2:   outputs["chart_version2"].Value,
		chartValuesYaml: outputs["chart_values_yaml"].Value,
		chartValuesJson: outputs["chart_values_json"].Value,
		ecsID:           outputs["ecs_id"].Value,
		ecsMirrorID:     outputs["ecs_mirror_id"].Value,
		ebmID:           outputs["ebm_id"].Value,
		ebmMirrorID:     outputs["ebm_mirror_id"].Value,
		ebmMirrorName:   outputs["ebm_mirror_name"].Value,
		deviceType:      outputs["device_type"].Value,
		ebmAz:           outputs["ebm_az"].Value,
	}
	fmt.Println("依赖资源初始化完毕")

	// 执行测试用例
	code := m.Run()

	// ccse依赖的子网无法马上删除, 所以不判断错误
	fmt.Println("开始清理依赖资源")
	// 清理依赖资源
	terraform.DestroyResource(dependenceDir)
	fmt.Println("依赖资源清理完毕")

	os.Exit(code)
}
