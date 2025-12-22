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
  name        = "tf-vpc-for-peer_connect"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_vpc" "vpc_test1" {
  name        = "tf-vpc-for-peer_connect1"
  cidr        = "172.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

# 不跨帐号
resource "ctyun_vpc_peer_connection" "example" {
  project_id     = "0"
  name           = "vpc-peer-conn-example"
  description    = "对等连接样例"
  request_vpc_id = ctyun_vpc.vpc_test.id
  accept_vpc_id  = ctyun_vpc.vpc_test1.id
}

data "ctyun_vpc_route_tables" "route_table_test" {
  vpc_id = ctyun_vpc.vpc_test.id
}

resource "ctyun_vpc_route_table_rule" "example" {
  ip_version     = 4
  next_hop_id    = ctyun_vpc_peer_connection.example.id
  destination    = "172.13.1.1/32"
  next_hop_type  = "vpcpeering"
  route_table_id = data.ctyun_vpc_route_tables.route_table_test.route_tables[0].route_table_id
  description    = "对等连接路由规则样例"
}
