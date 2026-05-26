package cloud_assistant_test

import (
	"fmt"
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunCloudAssistantInvocation_withExistingCommand(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	cmdRnd := utils.GenerateRandomString()
	cmdResourceName := "ctyun_cloud_assistant_command." + cmdRnd
	invocationResourceName := "ctyun_cloud_assistant_invocation." + rnd
	cmdFile := "resource_ctyun_cloud_assistant_command.tf"
	invocationFile := "resource_ctyun_cloud_assistant_invocation_with_command.tf"
	ecsID1 := dependence.ecsID1

	commandName := "tf-invoke-cmd-" + cmdRnd

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Step 1: 创建命令
			{
				Config: utils.LoadTestCase(cmdFile,
					cmdRnd,
					commandName,
					"shell",
					"echo hello world",
					"test command for invocation",
					60,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(cmdResourceName, "id"),
				),
			},
			// Step 2: 使用已有命令触发执行
			{
				Config: utils.LoadTestCase(cmdFile,
					cmdRnd,
					commandName,
					"shell",
					"echo hello world",
					"test command for invocation",
					60,
				) + utils.LoadTestCase(invocationFile,
					rnd,
					ecsID1,
					fmt.Sprintf("ctyun_cloud_assistant_command.%s.id", cmdRnd),
					60,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(invocationResourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(cmdFile,
					cmdRnd,
					commandName,
					"shell",
					"echo hello world",
					"test command for invocation",
					60,
				) + utils.LoadTestCase(invocationFile,
					rnd,
					ecsID1,
					fmt.Sprintf("ctyun_cloud_assistant_command.%s.id", cmdRnd),
					60,
				),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunCloudAssistantInvocation_withResultsDatasource(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resultsDsName := "data.ctyun_cloud_assistant_invocation_results." + dnd
	invocationDsFile := "resource_ctyun_cloud_assistant_invocation_with_results_ds.tf"
	ecsID := dependence.ecsID2
	commandID := dependence.commandID

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Step 1: 内联创建并触发命令 + 查询结果
			{
				Config: utils.LoadTestCase(invocationDsFile,
					rnd,
					ecsID,
					commandID,
					"/tmp",
					100,
					dnd,
					1,
					10,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resultsDsName, "total_count"),
					resource.TestCheckResourceAttrSet(resultsDsName, "results.#"),
				),
			},
		},
	})
}

func TestAccCtyunCloudAssistantInvocation_withParameterAndWorkingDir(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	invocationResourceName := "ctyun_cloud_assistant_invocation." + rnd
	invocationFile := "resource_ctyun_cloud_assistant_invocation_with_params.tf"

	cmdRnd := utils.GenerateRandomString()
	cmdResourceName := "ctyun_cloud_assistant_command." + cmdRnd
	cmdFile := "resource_ctyun_cloud_assistant_command_with_params.tf"

	ecsID1 := dependence.ecsID1
	ecsID2 := dependence.ecsID2
	ecsID := fmt.Sprintf("%s,%s", ecsID1, ecsID2)

	param := `{"name":"myapp","version":"v1"}`
	updatedParam := `{"name":"myapp","version":"v2"}`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// step1: 创建命令
			{
				Config: utils.LoadTestCase(cmdFile, cmdRnd, "tf-params-"+rnd, "shell", "echo {{name}}={{version}}", "带参数的触发测试", 30, true, param),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(cmdResourceName, "id"),
				),
			},
			// Step 1: 创建带参数和工作目录的命令执行
			{
				Config: utils.LoadTestCase(cmdFile, cmdRnd, "tf-params-"+rnd, "shell", "echo {{name}}={{version}}", "带参数的触发测试", 30, true, param) +
					utils.LoadTestCase(invocationFile,
						rnd,
						ecsID,
						fmt.Sprintf("%s.id", cmdResourceName),
						"/tmp",
						30,
						param,
					),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(invocationResourceName, "id"),
				),
			},
			// Step 2: 更新参数重新触发（RequiresReplace）
			{
				Config: utils.LoadTestCase(cmdFile, cmdRnd, "tf-params-"+rnd, "shell", "echo {{name}}={{version}}", "带参数的触发测试", 30, true, param) +
					utils.LoadTestCase(invocationFile,
						rnd,
						ecsID,
						fmt.Sprintf("%s.id", cmdResourceName),
						"/tmp",
						40,
						updatedParam,
					),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(invocationResourceName, "id"),
				),
			},
		},
	})
}
