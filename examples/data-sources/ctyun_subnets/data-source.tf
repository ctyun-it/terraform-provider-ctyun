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

resource "ctyun_vpc" "vpc_test" {
  name        = "vpc-test-mc22"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
  region_id   = "200000001852"
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "subnet-test"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
  enable_ipv6 = true
  region_id   = "200000001852"
}

data "ctyun_subnets" "test" {
  region_id = "200000001852"
  vpc_id = ctyun_vpc.vpc_test.id
  # page_no = 1
  # page_size = 1
}

output "ctyun_test" {
  value = data.ctyun_subnets.test
}

