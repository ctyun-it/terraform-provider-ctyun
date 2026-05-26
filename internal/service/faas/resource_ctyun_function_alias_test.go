package faas_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunFunctionAlias_Basic(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	rnd2 := utils.GenerateRandomString()
	resourceName := "ctyun_function_alias." + rnd
	versionId := "1" // 默认使用版本 1

	resourceFile := "resource_ctyun_function_alias.tf"
	aliasName := "alias-" + rnd
	description := "test-alias-" + rnd
	descriptionUpdate := "test-alias-" + rnd + "-updated"
	resourceFile_function := "resource_ctyun_function.tf"
	resourceName_function := "ctyun_function." + rnd
	resourceFile_version := "resource_ctyun_function_version.tf"
	resourceName_version := "ctyun_function_version." + rnd

	functionName := "func-" + utils.GenerateRandomString()
	runtimeRuntime := "python3.9"
	runtimeHandleType := "http"
	runtimeHandler := "index.handler"
	runtimeExecuteTimeout := 30
	runtimeInstanceConcurrency := 10
	containerTimeZone := "Asia/Shanghai"
	containerDiskSize := 512
	containerMemorySize := 2048
	containerCpu := 1.0
	containerListenPort := 8080
	codeBucket := "bucket-for-faas"
	codeKey := "hello_server.zip"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				// Step 1: 创建 Node.js 函数
				Config: utils.LoadTestCase(resourceFile_function,
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
					resource.TestCheckResourceAttr(resourceName_function, "name", functionName),
				),
			},
			{
				// Step 1: 创建 Node.js 函数
				Config: utils.LoadTestCase(resourceFile_function,
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
				) + utils.LoadTestCase(resourceFile_version,
					rnd,
					functionName,
					description,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName_function, "name", functionName),
				),
			},
			{
				// Step 2: 更新函数配置
				Config: utils.LoadTestCase(resourceFile_version,
					rnd,
					functionName,
					description,
				) + utils.LoadTestCase(resourceFile_function,
					rnd,
					functionName,
					runtimeRuntime,
					runtimeHandleType,
					runtimeHandler,
					runtimeExecuteTimeout,
					runtimeInstanceConcurrency,
					containerTimeZone,
					512, // 更新磁盘大小
					128, // 更新内存
					0.1, // 更新 CPU
					containerListenPort,
					description,
					codeBucket,
					codeKey,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName_function, "name", functionName),
				),
			}, {
				Config: utils.LoadTestCase(resourceFile_version,
					rnd,
					functionName,
					description,
				) + utils.LoadTestCase(resourceFile_function,
					rnd,
					functionName,
					runtimeRuntime,
					runtimeHandleType,
					runtimeHandler,
					runtimeExecuteTimeout,
					runtimeInstanceConcurrency,
					containerTimeZone,
					512, // 更新磁盘大小
					128, // 更新内存
					0.1, // 更新 CPU
					containerListenPort,
					description,
					codeBucket,
					codeKey,
				) + utils.LoadTestCase(resourceFile_version,
					rnd2,
					functionName,
					description+"2",
				),
				Check: resource.ComposeTestCheckFunc(),
			},
			{
				Config: utils.LoadTestCase(resourceFile_version,
					rnd,
					functionName,
					description,
				) + utils.LoadTestCase(resourceFile_function,
					rnd,
					functionName,
					runtimeRuntime,
					runtimeHandleType,
					runtimeHandler,
					runtimeExecuteTimeout,
					runtimeInstanceConcurrency,
					containerTimeZone,
					512, // 更新磁盘大小
					128, // 更新内存
					0.1, // 更新 CPU
					containerListenPort,
					description,
					codeBucket,
					codeKey,
				) + utils.LoadTestCase(resourceFile_version,
					rnd2,
					functionName,
					description+"2",
				) + utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					aliasName,
					versionId,
					description,
					resourceName_version,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "function_name", functionName),
					resource.TestCheckResourceAttr(resourceName, "alias_name", aliasName),
					resource.TestCheckResourceAttr(resourceName, "version_id", versionId),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile_version,
					rnd,
					functionName,
					description,
				) + utils.LoadTestCase(resourceFile_function,
					rnd,
					functionName,
					runtimeRuntime,
					runtimeHandleType,
					runtimeHandler,
					runtimeExecuteTimeout,
					runtimeInstanceConcurrency,
					containerTimeZone,
					512, // 更新磁盘大小
					128, // 更新内存
					0.1, // 更新 CPU
					containerListenPort,
					description,
					codeBucket,
					codeKey,
				) + utils.LoadTestCase(resourceFile_version,
					rnd2,
					functionName,
					description+"2",
				) + utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					aliasName,
					versionId,
					descriptionUpdate,
					resourceName_version,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdate),
				),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					function_name := ds.Attributes["function_name"]
					alias_name := ds.Attributes["alias_name"]
					return fmt.Sprintf("%s,%s", alias_name, function_name), nil
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile_version,
					rnd,
					functionName,
					description,
				) + utils.LoadTestCase(resourceFile_function,
					rnd,
					functionName,
					runtimeRuntime,
					runtimeHandleType,
					runtimeHandler,
					runtimeExecuteTimeout,
					runtimeInstanceConcurrency,
					containerTimeZone,
					512, // 更新磁盘大小
					128, // 更新内存
					0.1, // 更新 CPU
					containerListenPort,
					description,
					codeBucket,
					codeKey,
				) + utils.LoadTestCase(resourceFile_version,
					rnd2,
					functionName,
					description+"2",
				) + utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					aliasName,
					versionId,
					descriptionUpdate,
					resourceName_version,
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunFunctionAlias_WithGray(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	rnd2 := utils.GenerateRandomString()
	resourceName := "ctyun_function_alias." + rnd

	versionId := "1"
	grayVersionId := "2" // 灰度版本

	resourceFile := "resource_ctyun_function_alias_with_gray.tf"
	aliasName := "alias-gray-" + rnd
	description := "test-alias-" + rnd
	descriptionUpdate := "test-alias-" + rnd + "-updated"
	resourceFile_function := "resource_ctyun_function.tf"
	resourceName_function := "ctyun_function." + rnd
	resourceFile_version := "resource_ctyun_function_version.tf"
	resourceName_version := "ctyun_function_version." + rnd + "," + "ctyun_function_version." + rnd2

	functionName := "func-" + utils.GenerateRandomString()
	runtimeRuntime := "python3.9"
	runtimeHandleType := "http"
	runtimeHandler := "index.handler"
	runtimeExecuteTimeout := 30
	runtimeInstanceConcurrency := 10
	containerTimeZone := "Asia/Shanghai"
	containerDiskSize := 512
	containerMemorySize := 2048
	containerCpu := 1.0
	containerListenPort := 8080
	codeBucket := "bucket-for-faas"
	codeKey := "hello_server.zip"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				// Step 1: 创建 Node.js 函数
				Config: utils.LoadTestCase(resourceFile_function,
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
					resource.TestCheckResourceAttr(resourceName_function, "name", functionName),
				),
			},
			{
				// Step 1: 创建 Node.js 函数
				Config: utils.LoadTestCase(resourceFile_function,
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
				) + utils.LoadTestCase(resourceFile_version,
					rnd,
					functionName,
					description,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName_function, "name", functionName),
				),
			},
			{
				// Step 2: 更新函数配置
				Config: utils.LoadTestCase(resourceFile_version,
					rnd,
					functionName,
					description,
				) + utils.LoadTestCase(resourceFile_function,
					rnd,
					functionName,
					runtimeRuntime,
					runtimeHandleType,
					runtimeHandler,
					runtimeExecuteTimeout,
					runtimeInstanceConcurrency,
					containerTimeZone,
					512, // 更新磁盘大小
					128, // 更新内存
					0.1, // 更新 CPU
					containerListenPort,
					description,
					codeBucket,
					codeKey,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName_function, "name", functionName),
				),
			}, {
				Config: utils.LoadTestCase(resourceFile_version,
					rnd,
					functionName,
					description,
				) + utils.LoadTestCase(resourceFile_function,
					rnd,
					functionName,
					runtimeRuntime,
					runtimeHandleType,
					runtimeHandler,
					runtimeExecuteTimeout,
					runtimeInstanceConcurrency,
					containerTimeZone,
					512, // 更新磁盘大小
					128, // 更新内存
					0.1, // 更新 CPU
					containerListenPort,
					description,
					codeBucket,
					codeKey,
				) + utils.LoadTestCase(resourceFile_version,
					rnd2,
					functionName,
					description+"2",
				),
				Check: resource.ComposeTestCheckFunc(),
			},
			{
				Config: utils.LoadTestCase(resourceFile_version,
					rnd,
					functionName,
					description,
				) + utils.LoadTestCase(resourceFile_function,
					rnd,
					functionName,
					runtimeRuntime,
					runtimeHandleType,
					runtimeHandler,
					runtimeExecuteTimeout,
					runtimeInstanceConcurrency,
					containerTimeZone,
					512, // 更新磁盘大小
					128, // 更新内存
					0.1, // 更新 CPU
					containerListenPort,
					description,
					codeBucket,
					codeKey,
				) + utils.LoadTestCase(resourceFile_version,
					rnd2,
					functionName,
					description+"2",
				) + utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					versionId,
					grayVersionId,
					description,
					resourceName_version,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "function_name", functionName),
					resource.TestCheckResourceAttr(resourceName, "alias_name", aliasName),
					resource.TestCheckResourceAttr(resourceName, "version_id", versionId),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile_version,
					rnd,
					functionName,
					description,
				) + utils.LoadTestCase(resourceFile_function,
					rnd,
					functionName,
					runtimeRuntime,
					runtimeHandleType,
					runtimeHandler,
					runtimeExecuteTimeout,
					runtimeInstanceConcurrency,
					containerTimeZone,
					512, // 更新磁盘大小
					128, // 更新内存
					0.1, // 更新 CPU
					containerListenPort,
					description,
					codeBucket,
					codeKey,
				) + utils.LoadTestCase(resourceFile_version,
					rnd2,
					functionName,
					description+"2",
				) + utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					versionId,
					grayVersionId,
					description,
					resourceName_version,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "function_name", functionName),
					resource.TestCheckResourceAttr(resourceName, "alias_name", aliasName),
					resource.TestCheckResourceAttr(resourceName, "version_id", versionId),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "gray_version_id", grayVersionId),
					resource.TestCheckResourceAttr(resourceName, "gray_weight", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},

			{
				Config: utils.LoadTestCase(resourceFile_version,
					rnd,
					functionName,
					description,
				) + utils.LoadTestCase(resourceFile_function,
					rnd,
					functionName,
					runtimeRuntime,
					runtimeHandleType,
					runtimeHandler,
					runtimeExecuteTimeout,
					runtimeInstanceConcurrency,
					containerTimeZone,
					512, // 更新磁盘大小
					128, // 更新内存
					0.1, // 更新 CPU
					containerListenPort,
					description,
					codeBucket,
					codeKey,
				) + utils.LoadTestCase(resourceFile_version,
					rnd2,
					functionName,
					description+"2",
				) + utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					versionId,
					grayVersionId,
					descriptionUpdate,
					resourceName_version,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdate),
					resource.TestCheckResourceAttr(resourceName, "gray_weight", "10"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					function_name := ds.Attributes["function_name"]
					alias_name := ds.Attributes["alias_name"]
					return fmt.Sprintf("%s,%s", alias_name, function_name), nil
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile_version,
					rnd,
					functionName,
					description,
				) + utils.LoadTestCase(resourceFile_function,
					rnd,
					functionName,
					runtimeRuntime,
					runtimeHandleType,
					runtimeHandler,
					runtimeExecuteTimeout,
					runtimeInstanceConcurrency,
					containerTimeZone,
					512, // 更新磁盘大小
					128, // 更新内存
					0.1, // 更新 CPU
					containerListenPort,
					description,
					codeBucket,
					codeKey,
				) + utils.LoadTestCase(resourceFile_version,
					rnd2,
					functionName,
					description+"2",
				) + utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					versionId,
					grayVersionId,
					description,
					resourceName_version,
				),
				Destroy: true,
			},
		},
	})
}
