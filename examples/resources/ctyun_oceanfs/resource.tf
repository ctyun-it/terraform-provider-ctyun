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
  name        = "tf-vpc-for-sfs1"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}


resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-sfs"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}


resource "ctyun_oceanfs" "example" {
  protocol    = "nfs"
  name        = "oceanfs-examples"
  size        = "100"
  cycle_type  = "month"
  cycle_count = "1"
  vpc_id      = ctyun_vpc.vpc_test.id
  subnet_id   = ctyun_subnet.subnet_test.id
  tags        = [{ "key" : "test1", "value" : "value1" }, { "key" : "test2", "value" : "value2" }]
}