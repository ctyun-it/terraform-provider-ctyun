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

resource "ctyun_rabbitmq_exchange" "test" {
  instance_id = "8ccc8af2e6704080a72548735a081660"
  vhost       = "/"
  name        = "edxca"
  type        = "direct"
  # x_delayed_type = "fanout"
  # durable = true
  # auto_delete = true
  # internal = true
}