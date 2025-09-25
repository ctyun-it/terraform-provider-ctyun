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