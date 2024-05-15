terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
}

resource "ctyun_vpc" "vpc_test" {
  name        = "vpc-test-mc1"
  cidr        = "10.0.0.0/8"
  description = "terraform测试使用"
  enable_ipv6 = true
  project_id  = "4f5ef15300724760af59b37cf6409f45"
  region_id   = "200000002527"
}