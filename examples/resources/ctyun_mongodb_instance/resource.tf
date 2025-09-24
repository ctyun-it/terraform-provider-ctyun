terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 可参考index.md，在环境变量中配置ak、sk、资源池ID、可用区名称
provider "ctyun" {

}

resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-mon"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-mon"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-mon"
  description = "terraform测试使用"
  lifecycle {
    prevent_destroy = false
  }
}

resource "ctyun_mongodb_instance" "test" {
  cycle_type             = "on_demand"
  vpc_id                 = ctyun_vpc.vpc_test.id
  flavor_name            = "s7.large.2"
  subnet_id              = ctyun_subnet.subnet_test.id
  security_group_id      =  ctyun_security_group.security_group_test.id
  name                   = "mongodb-12ab"
  prod_id                = "Single34"
  storage_type           = "SATA"
  storage_space          = 100
  backup_storage_type    = "OS"
  password = var.password
}

variable "password" {
  type      = string
  sensitive = true
}