terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = "vpc-d7zxz8j05c"
  name        = "subnet-test"
  cidr        = "10.0.0.0/8"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
  enable_ipv6 = false
#  type        = "common"
  region_id   = "200000002527"
  project_id  = "4f5ef15300724760af59b37cf6409f45"
}