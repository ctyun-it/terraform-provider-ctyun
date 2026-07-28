package cloud_assistant_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunCloudAssistantRunCommand_basic(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_cloud_assistant_run_command." + rnd
	resourceFile := "resource_ctyun_cloud_assistant_run_command.tf"
	ecsID := dependence.ecsID1

	commandName := "tf-run-cmd-" + rnd

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Step 1: 执行Shell命令（不保存）
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					ecsID,
					commandName,
					"shell",
					"echo hello from run_command",
					"terraform run_command test",
					false,
					60,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
				),
			},
		},
	})
}

func TestAccCtyunCloudAssistantRunCommand_saveCommand(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_cloud_assistant_run_command." + rnd
	resourceFile := "resource_ctyun_cloud_assistant_run_command.tf"
	ecsID := dependence.ecsID1

	commandName := "tf-save-cmd-" + rnd

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 执行命令并保存
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					ecsID,
					commandName,
					"shell",
					"echo save command test",
					"terraform save_command test",
					true,
					60,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "save_command", "true"),
				),
			},
		},
	})
}

func TestAccCtyunCloudAssistantRunCommand_allTypes(t *testing.T) {
	t.Parallel()
	typeCases := map[string]string{
		"shell":      "echo hello from shell",
		"bat":        "@echo hello from bat",
		"powershell": `Write-Host \"hello from powershell\"`,
		"python":     `print(\"hello from python\")`,
	}
	for cmdType, content := range typeCases {
		t.Run(cmdType, func(t *testing.T) {
			t.Parallel()
			testRunCommandType(t, cmdType, content)
		})
	}
}

func testRunCommandType(t *testing.T, cmdType string, content string) {
	t.Helper()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_cloud_assistant_run_command." + rnd
	resourceFile := "resource_ctyun_cloud_assistant_run_command.tf"
	ecsID := dependence.ecsID1

	commandName := fmt.Sprintf("tf-run-%s-%s", strings.ToLower(cmdType), rnd)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					ecsID,
					commandName,
					cmdType,
					content,
					fmt.Sprintf("test %s type run_command", cmdType),
					false,
					60,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "command_type", cmdType),
				),
			},
		},
	})
}

func TestAccCtyunCloudAssistantRunCommand_withParameter(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_cloud_assistant_run_command." + rnd
	resourceFile := "resource_ctyun_cloud_assistant_run_command_with_params.tf"
	ecsID := dependence.ecsID1

	commandName := "tf-param-cmd-" + rnd
	param := `{"name":"myapp","version":"v1"}`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Step 1: 带参数和工作目录执行
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					ecsID,
					commandName,
					"shell",
					"echo {{name}}={{version}}",
					"run_command with parameter test",
					"/tmp",
					true,
					true,
					60,
					param,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "enabled_parameter", "true"),
					resource.TestCheckResourceAttr(resourceName, "working_directory", "/tmp"),
				),
			},
		},
	})
}

func TestAccCtyunCloudAssistantRunCommand_withResultsDatasource(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dsrnd := utils.GenerateRandomString()
	resourceName := "ctyun_cloud_assistant_run_command." + rnd
	resultsDsName := "data.ctyun_cloud_assistant_invocation_results." + dsrnd
	resourceFile := "resource_ctyun_cloud_assistant_run_command_with_ds.tf"
	ecsID := dependence.ecsID2

	commandName := "tf-ds-cmd-" + rnd

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 执行命令 + 查询结果
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					ecsID,
					commandName,
					"shell",
					"echo datasource test",
					"run_command with datasource test",
					true,
					60,
					dsrnd,
					1,
					10,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resultsDsName, "results.#"),
				),
			},
		},
	})
}

func TestAccCtyunCloudAssistantRunCommand_multiInstance(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_cloud_assistant_run_command." + rnd
	resourceFile := "resource_ctyun_cloud_assistant_run_command.tf"
	ecsID := fmt.Sprintf("%s,%s", dependence.ecsID1, dependence.ecsID2)

	commandName := "tf-multi-cmd-" + rnd

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 多台实例同时执行
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					ecsID,
					commandName,
					"shell",
					"hostname",
					"multi instance test",
					false,
					60,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}
