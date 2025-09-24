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

resource "ctyun_security_group" "security_group_test" {
  vpc_id      = "vpc-r7kv00qbz5"
  name        = "terraform-minchiang-test55"
  description = "terraform测试使用"
}