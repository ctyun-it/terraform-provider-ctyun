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

variable "password" {
  type      = string
  sensitive = true
}

resource "ctyun_redis_account" "test" {
  name        = "test"
  instance_id = "425c9173f98b4646a72ce0b986af00b0"
  password    = var.password
  privilege   = "ro"
  description = "description1111"
}