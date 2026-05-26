terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 可参考index.md，在环境变量中配置ak、sk、资源池ID、可用区名称
provider "ctyun" {
  env = "prod"
}

# 管理云助手命令
resource "ctyun_cloud_assistant_command" "test" {
  command_name    = "cmd_example"
  command_type    = "shell"
  command_content = "echo hello world"
  description     = "云助手命令terraform样例"
  timeout         = 30
}
