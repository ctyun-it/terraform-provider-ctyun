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

data "ctyun_rabbitmq_instances" "tbidgqvfbs" {
  instance_id ="8d839e64a4314edb8121d0d1f69b8b19"
}

output "list" {
  value = data.ctyun_rabbitmq_instances.tbidgqvfbs
}