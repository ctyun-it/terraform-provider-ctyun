resource "ctyun_vpc" "vpc_test" {
  name        = "vpc-test-route-2"
  cidr        = "192.168.4.0/24"
  description = "terraform测试使用"
  enable_ipv6 = true
  lifecycle {
    create_before_destroy = true
  }
}

resource "ctyun_vpc_route_table" "route" {
  vpc_id = ctyun_vpc.vpc_test.id
  name = "route-tf-123"
}

resource "ctyun_nat" "nat_test2"{
  vpc_id = ctyun_vpc.vpc_test.id
  spec = 1
  name = "tf-nat2-route"
  description = "terraform测试使用"
  cycle_type = "on_demand"
}

resource "ctyun_vpc_route_table_rule" "rule_cn2-public-service2"{
  count = length(var.additional_routes)
  description = var.additional_routes[count.index].description
  destination = var.additional_routes[count.index].cidr_block
  next_hop_id = ctyun_nat.nat_test2.id
  next_hop_type = "natgw"
  route_table_id = ctyun_vpc_route_table.route.id
  ip_version = 4
}
