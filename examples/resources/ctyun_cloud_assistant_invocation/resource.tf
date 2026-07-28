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

resource "ctyun_cloud_assistant_command" "test" {
  command_name    = "cmd_example"
  command_type    = "shell"
  command_content = "echo hello world"
  description     = "云助手命令terraform样例"
  timeout         = 30
}


data "ctyun_ecs_instances" "test" {
  page_size = 1
}

output "ctyun_ecs_instances_test" {
  value = data.ctyun_ecs_instances.test
}

# 触发云助手命令执行
resource "ctyun_cloud_assistant_invocation" "test" {
  instance_ids    = data.ctyun_ecs_instances.test.instances[0].id
  command_id      = ctyun_cloud_assistant_command.test.id
  timeout         = 60
}
