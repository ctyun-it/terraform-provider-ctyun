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
  name        = "tf-vpc-for-nat"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_nat" "nat_test"{
  vpc_id = ctyun_vpc.vpc_test.id
  spec = 1
  name = "tf-nat"
  description = "terraform测试使用"
  cycle_type = "on_demand"
}

resource "ctyun_eip" "eip_test" {
  name                = "tf-eip-for-nat1"
  bandwidth           = 1
  cycle_type          = "on_demand"
  demand_billing_type = "upflowc"
}

resource "ctyun_nat" "nat_test"{
  vpc_id = ctyun_vpc.vpc_test.id
  spec = 1
  name = "tf-nat-expample"
  description = "terraform测试使用"
  cycle_type = "on_demand"
}

resource "ctyun_nat_dnat" "dnat_test"{
  nat_gateway_id = ctyun_nat.nat_test.id
  external_id = ctyun_eip.eip_test.id
  external_port = 80
  internal_ip = "127.0.0.1"
  dnat_type = "custom"
  internal_port = 12454
  protocol = "tcp"
}