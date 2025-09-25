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
  name        = "tf-vpc-for-elb"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-elb"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_elb_loadbalancer" "elb_test" {
  subnet_id     = ctyun_subnet.subnet_test.id
  name          = "tf-elb-for-test"
  sla_name      = "elb.s2.small"
  resource_type = "internal"
  vpc_id        = ctyun_vpc.vpc_test.id
  cycle_type    = "month"
  cycle_count   = 1
}
