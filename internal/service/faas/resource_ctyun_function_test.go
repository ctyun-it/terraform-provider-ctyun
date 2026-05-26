package faas_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunFunction_Python(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_function." + rnd
	dataSourceName := "data.ctyun_functions." + rnd

	resourceFile := "resource_ctyun_function.tf"
	dataSourceFile := "datasource_ctyun_functions.tf"

	functionName := "func-" + utils.GenerateRandomString()
	runtimeRuntime := "python3.9"
	runtimeHandleType := "event"
	runtimeHandler := "index.handler"
	runtimeExecuteTimeout := 60
	runtimeInstanceConcurrency := 1
	containerTimeZone := "UTC"
	containerDiskSize := 512
	containerMemorySize := 256
	containerCpu := 0.2
	containerListenPort := 8080
	description := "Test function created by Terraform"
	codeBucket := "bucket-faas-test"
	codeKey := "utils.zip"

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
			{
				// Step 1: 创建函数
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					runtimeRuntime,
					runtimeHandleType,
					runtimeHandler,
					runtimeExecuteTimeout,
					runtimeInstanceConcurrency,
					containerTimeZone,
					containerDiskSize,
					containerMemorySize,
					containerCpu,
					containerListenPort,
					description,
					codeBucket,
					codeKey,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", functionName),
					resource.TestCheckResourceAttr(resourceName, "runtime_runtime", runtimeRuntime),
					resource.TestCheckResourceAttr(resourceName, "runtime_handle_type", runtimeHandleType),
					resource.TestCheckResourceAttr(resourceName, "runtime_handler", runtimeHandler),
					resource.TestCheckResourceAttr(resourceName, "runtime_execute_timeout", "60"),
					resource.TestCheckResourceAttr(resourceName, "runtime_instance_concurrency", "1"),
					resource.TestCheckResourceAttr(resourceName, "container_time_zone", containerTimeZone),
					resource.TestCheckResourceAttr(resourceName, "container_disk_size", "512"),
					resource.TestCheckResourceAttr(resourceName, "container_memory_size", "256"),
					resource.TestCheckResourceAttr(resourceName, "container_cpu", "0.2"),
					resource.TestCheckResourceAttr(resourceName, "container_listen_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "function_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				// Step 2: 更新函数配置
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					runtimeRuntime,
					runtimeHandleType,
					runtimeHandler,
					runtimeExecuteTimeout,
					runtimeInstanceConcurrency,
					containerTimeZone,
					10240, // 更新磁盘大小
					512,   // 更新内存
					0.5,   // 更新 CPU
					containerListenPort,
					description+" - updated",
					codeBucket,
					codeKey,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", functionName),
					resource.TestCheckResourceAttr(resourceName, "container_disk_size", "10240"),
					resource.TestCheckResourceAttr(resourceName, "container_memory_size", "512"),
					resource.TestCheckResourceAttr(resourceName, "container_cpu", "0.5"),
					resource.TestCheckResourceAttr(resourceName, "description", description+" - updated"),
				),
			},
			{
				// Step 3: 数据源查询测试
				Config: utils.LoadTestCase(dataSourceFile,
					rnd,
					functionName,
				) + "\n" + utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					runtimeRuntime,
					runtimeHandleType,
					runtimeHandler,
					runtimeExecuteTimeout,
					runtimeInstanceConcurrency,
					containerTimeZone,
					10240,
					512,
					0.5,
					containerListenPort,
					description+" - updated",
					codeBucket,
					codeKey,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "functions.#", "1"),
				),
			},
			{
				// Step 4: 导入测试
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code_bucket", "code_key"},
			},

			{
				Config:  utils.LoadTestCase(resourceFile, rnd, functionName, runtimeRuntime, runtimeHandleType, runtimeHandler, runtimeExecuteTimeout, runtimeInstanceConcurrency, containerTimeZone, containerDiskSize, containerMemorySize, containerCpu, containerListenPort, description, codeBucket, codeKey),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunFunction_Nodejs(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_function." + rnd
	dataSourceName := "data.ctyun_functions." + rnd

	resourceFile := "resource_ctyun_function.tf"
	dataSourceFile := "datasource_ctyun_functions.tf"

	functionName := "func-" + utils.GenerateRandomString()
	runtimeRuntime := "nodejs16"
	runtimeHandleType := "http"
	runtimeHandler := "index.handler"
	runtimeExecuteTimeout := 30
	runtimeInstanceConcurrency := 10
	containerTimeZone := "Asia/Shanghai"
	containerDiskSize := 512
	containerMemorySize := 2048
	containerCpu := 1.0
	containerListenPort := 8080
	description := "Node.js test function"
	codeBucket := "bucket-faas-test"
	codeKey := "utils.zip"

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
			{
				// Step 1: 创建 Node.js 函数
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					runtimeRuntime,
					runtimeHandleType,
					runtimeHandler,
					runtimeExecuteTimeout,
					runtimeInstanceConcurrency,
					containerTimeZone,
					containerDiskSize,
					containerMemorySize,
					containerCpu,
					containerListenPort,
					description,
					codeBucket,
					codeKey,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", functionName),
					resource.TestCheckResourceAttr(resourceName, "runtime_runtime", runtimeRuntime),
					resource.TestCheckResourceAttr(resourceName, "runtime_handle_type", runtimeHandleType),
					resource.TestCheckResourceAttr(resourceName, "runtime_handler", runtimeHandler),
					resource.TestCheckResourceAttr(resourceName, "runtime_execute_timeout", "30"),
					resource.TestCheckResourceAttr(resourceName, "runtime_instance_concurrency", "10"),
					resource.TestCheckResourceAttr(resourceName, "container_time_zone", containerTimeZone),
					resource.TestCheckResourceAttr(resourceName, "container_disk_size", "512"),
					resource.TestCheckResourceAttr(resourceName, "container_memory_size", "2048"),
					resource.TestCheckResourceAttr(resourceName, "container_cpu", "1"),
					resource.TestCheckResourceAttr(resourceName, "container_listen_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "function_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				// Step 2: 更新函数配置
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					runtimeRuntime,
					runtimeHandleType,
					runtimeHandler,
					runtimeExecuteTimeout,
					runtimeInstanceConcurrency,
					containerTimeZone,
					10240, // 更新磁盘大小
					8064,  // 更新内存
					2.0,   // 更新 CPU
					containerListenPort,
					description+" - updated",
					codeBucket,
					codeKey,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", functionName),
					resource.TestCheckResourceAttr(resourceName, "container_disk_size", "10240"),
					resource.TestCheckResourceAttr(resourceName, "container_memory_size", "8064"),
					resource.TestCheckResourceAttr(resourceName, "container_cpu", "2"),
					resource.TestCheckResourceAttr(resourceName, "description", description+" - updated"),
				),
			},
			{
				// Step 3: 数据源查询测试
				Config: utils.LoadTestCase(dataSourceFile,
					rnd,

					functionName,
				) + "\n" + utils.LoadTestCase(resourceFile,
					rnd,
					functionName,

					runtimeRuntime,
					runtimeHandleType,
					runtimeHandler,
					runtimeExecuteTimeout,
					runtimeInstanceConcurrency,
					containerTimeZone,
					10240, // 更新磁盘大小
					8064,  // 更新内存
					2.0,   // 更新 CPU
					containerListenPort,
					description+" - updated",
					codeBucket,
					codeKey,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "functions.#", "1"),
				),
			},
			{
				// Step 4: 导入测试
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code_bucket", "code_key"},
			},
			{
				// Step 5: 销毁测试
				Config:  utils.LoadTestCase(resourceFile, rnd, functionName, runtimeRuntime, runtimeHandleType, runtimeHandler, runtimeExecuteTimeout, runtimeInstanceConcurrency, containerTimeZone, containerDiskSize, containerMemorySize, containerCpu, containerListenPort, description, codeBucket, codeKey),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunFunction_EnvironmentVariables(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_function." + rnd

	resourceFile := "resource_ctyun_function_env.tf"

	functionName := "func-env-" + utils.GenerateRandomString()
	runtimeRuntime := "python3.9"
	runtimeHandleType := "event"
	runtimeHandler := "index.handler"
	runtimeExecuteTimeout := 60
	runtimeInstanceConcurrency := 1
	containerTimeZone := "UTC"
	containerDiskSize := 512
	containerMemorySize := 256
	containerCpu := 0.2
	containerListenPort := 8080
	description := "Test function with environment variables"
	codeBucket := "bucket-faas-test"
	codeKey := "utils.zip"
	environment := `{"TEST_KEY":"test_value","API_KEY":"123456","LOG_LEVEL":"debug"}`

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
			{
				// Step 1: 创建函数 with EnvironmentVariables
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					runtimeRuntime,
					runtimeHandleType,
					runtimeHandler,
					runtimeExecuteTimeout,
					runtimeInstanceConcurrency,
					containerTimeZone,
					containerDiskSize,
					containerMemorySize,
					containerCpu,
					containerListenPort,
					description,
					codeBucket,
					codeKey,
					environment,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", functionName),
					resource.TestCheckResourceAttr(resourceName, "runtime_runtime", runtimeRuntime),
					resource.TestCheckResourceAttr(resourceName, "runtime_handler", runtimeHandler),
					resource.TestCheckResourceAttr(resourceName, "container_time_zone", containerTimeZone),
					resource.TestCheckResourceAttr(resourceName, "container_disk_size", "512"),
					resource.TestCheckResourceAttr(resourceName, "container_memory_size", "256"),
					resource.TestCheckResourceAttr(resourceName, "container_cpu", "0.2"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "environment.TEST_KEY", "test_value"),
					resource.TestCheckResourceAttr(resourceName, "environment.API_KEY", "123456"),
					resource.TestCheckResourceAttr(resourceName, "environment.LOG_LEVEL", "debug"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				// Step 2: 更新 EnvironmentVariables
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					runtimeRuntime,
					runtimeHandleType,
					runtimeHandler,
					runtimeExecuteTimeout,
					runtimeInstanceConcurrency,
					containerTimeZone,
					containerDiskSize,
					containerMemorySize,
					containerCpu,
					containerListenPort,
					description+" - updated",
					codeBucket,
					codeKey,
					`{"TEST_KEY":"updated_value","NEW_KEY":"new_value"}`,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "environment.TEST_KEY", "updated_value"),
					resource.TestCheckResourceAttr(resourceName, "environment.NEW_KEY", "new_value"),
					resource.TestCheckNoResourceAttr(resourceName, "environment.API_KEY"),
					resource.TestCheckNoResourceAttr(resourceName, "environment.LOG_LEVEL"),
				),
			},
			{
				// Step 3: 销毁测试
				Config:  utils.LoadTestCase(resourceFile, rnd, functionName, runtimeRuntime, runtimeHandleType, runtimeHandler, runtimeExecuteTimeout, runtimeInstanceConcurrency, containerTimeZone, containerDiskSize, containerMemorySize, containerCpu, containerListenPort, description, codeBucket, codeKey, environment),
				Destroy: true,
			},
		},
	})
}
