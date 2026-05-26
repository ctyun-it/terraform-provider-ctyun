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


data "ctyun_ecs_instances" "test" {
  page_size = 1
}

output "ctyun_ecs_instances_test" {
  value = data.ctyun_ecs_instances.test
}
# 执行云助手命令（Shell类型）
resource "ctyun_cloud_assistant_run_command" test {
  instance_ids    = data.ctyun_ecs_instances.test.instances[0].id
  command_name    = "hello_terraform"
  command_type    = "shell"
  command_content = "echo hello terraform"
  description     = "运行命令样例"
  save_command    = false
  timeout         = 30
}
