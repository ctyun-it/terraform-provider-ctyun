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
  name        = "vpc-test-mc1"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_vpc_route_table" "route" {
  vpc_id = ctyun_vpc.vpc_test.id
  name = "route-t1f"
}

data "ctyun_vpc_route_table_rules" "rtest" {
  route_table_id = ctyun_vpc_route_table.route.route_table_id
}

locals {
  igw_rules = [for rule in data.ctyun_vpc_route_table_rules.rtest.rules : rule if rule.next_hop_type == "igw"]

  igw_id = length(local.igw_rules) > 0 ? local.igw_rules[0].next_hop_id : ""
}

resource "ctyun_vpc_route_table_rule" "rule_test"{
  description = "test"
  destination = "188.168.0.0/16"
  next_hop_id = local.igw_id
  next_hop_type = "igw"
  route_table_id = ctyun_vpc_route_table.route.route_table_id
  ip_version = 4
}