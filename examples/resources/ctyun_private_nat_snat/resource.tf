terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

provider "ctyun" {
  env = "prod"
}
resource "ctyun_vpc" "vpc_test" {
  name        = "vpc-test-ccse1"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "subnet-test-ccse1"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns = [
    "100.95.0.1"
  ]
  enable_ipv6 = true
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-image"
  description = "terraform测试使用"
}
resource "ctyun_port" "port" {
  name                       = "port-test-update"
  description                = "port 测试-测试"
  subnet_id                  = ctyun_subnet.subnet_test.id
  security_group_ids         = [ctyun_security_group.security_group_test.id]
  secondary_private_ip_count = 1
}
resource "ctyun_private_nat" "private_nat" {
  vpc_id      = ctyun_vpc.vpc_test.id
  spec        = "small"
  name        = "private-nat-test"
  description = "私有网关测试"
  cycle_type  = "on_demand"
  subnet_id   = ctyun_subnet.subnet_test.id
}

resource "ctyun_private_nat_transit_ip" "transit_ip" {
  nat_gateway_id = ctyun_private_nat.private_nat.id
  address        = "192.168.0.5"
}
resource "ctyun_private_nat_transit_ip" "transit_ip2" {
  nat_gateway_id = ctyun_private_nat.private_nat.id
  address        = "192.168.0.8"
}

resource "ctyun_private_nat_snat" "private_nat_snat" {
  nat_gateway_id   = ctyun_private_nat.private_nat.id
  source_subnet_id = ctyun_subnet.subnet_test.id
  addresses        = [ctyun_private_nat_transit_ip.transit_ip.address, ctyun_private_nat_transit_ip.transit_ip2.address]
  description      = "snat 测试-更新"
}
