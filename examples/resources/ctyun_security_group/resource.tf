terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id      = "vpc-r7kv00qbz5"
  name        = "terraform-minchiang-test55"
  description = "terraform测试使用"
}