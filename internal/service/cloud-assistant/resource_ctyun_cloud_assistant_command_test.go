package cloud_assistant_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCtyunCloudAssistantCommand_basic(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()
	resourceName := "ctyun_cloud_assistant_command." + rnd
	resourceFile := "resource_ctyun_cloud_assistant_command.tf"

	datasourceName := "data.ctyun_cloud_assistant_commands." + dnd
	datasourceFile := "datasource_ctyun_cloud_assistant_commands.tf"

	commandName := "tf-test-cmd-" + rnd
	updatedName := "tf-test-cmd-upd-" + rnd
	commandType := "shell"
	commandContent := "echo hello"
	updatedCommandContent := "echo hello world"
	description := "terraform test command"
	updatedDescription := "terraform test updated"
	timeout := 60
	updatedTimeout := 80

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
			// 1. 创建命令
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					commandName,
					commandType,
					commandContent,
					description,
					timeout,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "command_name", commandName),
					resource.TestCheckResourceAttr(resourceName, "command_type", commandType),
					resource.TestCheckResourceAttr(resourceName, "command_content", commandContent),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "timeout", fmt.Sprintf("%d", timeout)),
				),
			},
			// 2. 导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s",
						rs.Primary.Attributes["id"],
						rs.Primary.Attributes["region_id"],
					), nil
				},
				ImportStateVerify: true,
			},
			// 3. 更新命令名称和描述
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					updatedName,
					commandType,
					updatedCommandContent,
					updatedDescription,
					updatedTimeout,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "command_name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
			},
			// 4. 验证datasource（查询命令列表）
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					updatedName,
					commandType,
					commandContent,
					updatedDescription,
					timeout,
				) + utils.LoadTestCase(datasourceFile,
					dnd,
					1,
					10,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "total_count"),
					resource.TestCheckResourceAttrSet(datasourceName, "commands.#"),
				),
			},
			// 5. 销毁
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, updatedName, commandType, updatedCommandContent, updatedDescription, updatedTimeout),
				Destroy: true,
			},
		},
	})
}

func TestAccCtyunCloudAssistantCommand_withParameter(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_cloud_assistant_command." + rnd
	resourceFile := "resource_ctyun_cloud_assistant_command_with_params.tf"

	commandName := "tf-test-param-" + rnd
	updatedName := "tf-test-param-upd-" + rnd
	commandType := "shell"
	commandContent := `echo {{name}} {{version}}`
	description := "param test"
	updatedDescription := "param test updated"
	timeout := 60
	parameter := fmt.Sprintf(`{"%s":"%s","%s":"%s","%s":"%s"}`, "name", "myapp", "version", "v1", "owner", "kqjwyk")
	updatedParameter := fmt.Sprintf(`{"%s":"%s","%s":"%s"}`, "name", "myapp", "version", "v2")

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
			// 1. 创建带参数的命令
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					commandName,
					commandType,
					commandContent,
					description,
					timeout,
					true,
					parameter,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "command_name", commandName),
					resource.TestCheckResourceAttr(resourceName, "enabled_parameter", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			// 2. 更新参数和描述
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					updatedName,
					commandType,
					commandContent,
					updatedDescription,
					timeout,
					true,
					updatedParameter,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "command_name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
			},

			// 4. 销毁
			{
				Destroy: true,
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					updatedName,
					commandType,
					commandContent,
					updatedDescription,
					timeout,
					true,
					updatedParameter,
				),
			},
		},
	})
}

func TestAccCtyunCloudAssistantCommand_allTypes(t *testing.T) {
	t.Parallel()
	typeCases := map[string]string{
		"shell":      "echo hello from shell",
		"bat":        "@echo hello from bat",
		"powershell": `Write-Host \"hello from powershell\"`,
		"python":     `print('hello from python')`,
	}
	for cmdType, content := range typeCases {
		t.Run(cmdType, func(t *testing.T) {
			testCommandType(t, cmdType, content)
		})
	}
}

func testCommandType(t *testing.T, cmdType string, content string) {
	t.Helper()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_cloud_assistant_command." + rnd
	resourceFile := "resource_ctyun_cloud_assistant_command.tf"

	commandName := fmt.Sprintf("tf-cmd-%s-%s", strings.ToLower(cmdType), rnd)

	resource.Test(t, resource.TestCase{
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed for type %s", cmdType)
			}
			return nil
		},
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 创建命令
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					commandName,
					cmdType,
					content,
					fmt.Sprintf("test %s type command", cmdType),
					60,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "command_name", commandName),
					resource.TestCheckResourceAttr(resourceName, "command_type", cmdType),
					//resource.TestCheckResourceAttr(resourceName, "command_content", content),
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("test %s type command", cmdType)),
				),
			},
			// 销毁
			{
				Destroy: true,
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					commandName,
					cmdType,
					content,
					fmt.Sprintf("test %s type command", cmdType),
					60,
				),
			},
		},
	})
}
