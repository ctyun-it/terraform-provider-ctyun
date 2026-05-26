package opensearch_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"os"
	"testing"
)

const dependenceDir = "testdata/dependence"

type Dependence struct {
	regionID        string
	vpcID           string
	subnetID        string
	securityGroupID string
	azName          string
	flavorName      string
	storageType     string
}

var dependence Dependence

func TestMain(m *testing.M) {
	if skip := os.Getenv("SKIP_OPENSEARCH_TEST"); skip != "" {
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
		regionID:        toString(outputs["region_id"].Value),
		vpcID:           toString(outputs["vpc_id"].Value),
		subnetID:        toString(outputs["subnet_id"].Value),
		securityGroupID: toString(outputs["security_group_id"].Value),
		azName:          toString(outputs["az_name"].Value),
		flavorName:      toString(outputs["flavor_name"].Value),
		storageType:     toString(outputs["storage_type"].Value),
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

// toString converts terraform output values to string
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case fmt.Stringer:
		return val.String()
	default:
		return fmt.Sprintf("%v", val)
	}
}
