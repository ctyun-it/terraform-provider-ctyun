package elb_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

func TestAccCtyunElbHealthCheck(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_elb_health_check." + rnd
	resourceFile := "resource_ctyun_elb_health_check.tf"

	datasourceName := "data.ctyun_elb_health_checks." + dnd
	datasourceFile := "datasource_ctyun_elb_health_checks.tf"

	name := "health_check_" + utils.GenerateRandomString()
	updatedName := name + "_new"
	protocol := "TCP"

	timeout := 60
	interval := 60
	maxRetry := 10
	httpMethod := "POST"
	httpUrlPath := "/health"
	httpExpectedCodes := `"http_2xx","http_3xx","http_4xx","http_5xx"`
	protocolPort := 8080

	timeoutTF := fmt.Sprintf(`timeout=%d`, timeout)
	intervalTF := fmt.Sprintf(`interval=%d`, interval)
	maxRetryTF := fmt.Sprintf(`max_retry=%d`, maxRetry)
	httpMethodTF := fmt.Sprintf(`http_method="%s"`, httpMethod)
	httpUrlPathTF := fmt.Sprintf(`http_url_path="%s"`, httpUrlPath)
	httpExpectedCodesTF := fmt.Sprintf(`http_expected_codes=[%s]`, httpExpectedCodes)
	protocolPortTF := fmt.Sprintf(`protocol_port=%d`, protocolPort)

	updatedTimeout := 59
	updatedInterval := 59
	updatedMaxRetry := 9
	updatedHttpMethod := "DELETE"
	updatedHttpUrlPath := "/health/test"
	updatedHttpExpectedCodes := `"http_2xx","http_3xx"`

	updatedTimeoutTF := fmt.Sprintf(`timeout=%d`, updatedTimeout)
	updatedIntervalTF := fmt.Sprintf(`interval=%d`, updatedInterval)
	updatedMaxRetryTF := fmt.Sprintf(`max_retry=%d`, updatedMaxRetry)
	updatedHttpMethodTF := fmt.Sprintf(`http_method="%s"`, updatedHttpMethod)
	updatedHttpUrlPathTF := fmt.Sprintf(`http_url_path="%s"`, updatedHttpUrlPath)
	updatedHttpExpectedCodesTF := fmt.Sprintf(`http_expected_codes=[%s]`, updatedHttpExpectedCodes)

	resource.Test(t, resource.TestCase{
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1 基础测试，不包含各类超时参数
			// 1.1验证创建健康检查
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name, protocol, "", "", "", "", "", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "protocol", protocol),
				),
			},
			// 1.2 验证更新健康检查
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, protocol, "", "", "", "", "", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "protocol", protocol),
				),
			},
			// 1.3 验证datasource
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, protocol, "", "", "", "", "", "", "") +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf("ids=%s.id", resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr(datasourceName, "health_checks.#", "1"),
					//resource.TestCheckResourceAttr(datasourceName, "health_checks.0.name", updatedName),
					resource.TestCheckResourceAttr(datasourceName, "health_checks.0.protocol", protocol),
				),
			},

			// 1.4 验证删除健康检查
			{
				Config:  utils.LoadTestCase(resourceFile, dnd, updatedName, protocol, "", "", "", "", "", "", ""),
				Destroy: true,
			},

			// 2protocol=http测试，再加上各类超时参数
			// 2.1 创建测试
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name, "HTTP", timeoutTF, intervalTF, maxRetryTF, httpMethodTF, httpUrlPathTF, httpExpectedCodesTF, protocolPortTF),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "timeout", strconv.Itoa(timeout)),
					resource.TestCheckResourceAttr(resourceName, "interval", strconv.Itoa(interval)),
					resource.TestCheckResourceAttr(resourceName, "max_retry", strconv.Itoa(maxRetry)),
					resource.TestCheckResourceAttr(resourceName, "http_method", httpMethod),
					resource.TestCheckResourceAttr(resourceName, "http_url_path", httpUrlPath),
					resource.TestCheckResourceAttr(resourceName, "http_expected_codes.#", strconv.Itoa(4)),
					resource.TestCheckResourceAttr(resourceName, "protocol_port", strconv.Itoa(protocolPort)),
				),
			},
			// 2.2 update 测试
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name, "HTTP", updatedTimeoutTF, updatedIntervalTF, updatedMaxRetryTF, updatedHttpMethodTF, updatedHttpUrlPathTF, updatedHttpExpectedCodesTF, protocolPortTF),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "timeout", strconv.Itoa(updatedTimeout)),
					resource.TestCheckResourceAttr(resourceName, "interval", strconv.Itoa(updatedInterval)),
					resource.TestCheckResourceAttr(resourceName, "max_retry", strconv.Itoa(updatedMaxRetry)),
					resource.TestCheckResourceAttr(resourceName, "http_method", updatedHttpMethod),
					resource.TestCheckResourceAttr(resourceName, "http_url_path", updatedHttpUrlPath),
					resource.TestCheckResourceAttr(resourceName, "http_expected_codes.#", strconv.Itoa(2)),
					resource.TestCheckResourceAttr(resourceName, "protocol_port", strconv.Itoa(protocolPort)),
				),
			},
			// 2.3 delete验证
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, name, "HTTP", updatedTimeoutTF, updatedIntervalTF, updatedMaxRetryTF, updatedHttpMethodTF, updatedHttpUrlPathTF, updatedHttpExpectedCodesTF, protocolPortTF),
				Destroy: true,
			},
		},
	})
}
