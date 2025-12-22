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

resource "ctyun_vpc" "vpc_test3" {
  name        = "tf-vpc-for-peer_connect2"
  cidr        = "172.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_vpc_peer_connection" "test" {
  project_id     = "0"
  name           = "peer-connection_tf"
  description    = "terraform测试使用"
  request_vpc_id = ctyun_vpc.vpc_test.id
  accept_vpc_id  = ctyun_vpc.vpc_test3.id
}


# provider "ctyun" {
#   alias           = "test_accpet"
#   ak              = "xxx"                                    # 如果此值不填，则默认读取环境变量中的CTYUN_AK
#   sk              = "xxx"                                    # 如果此值不填，则默认读取环境变量中的CTYUN_SK
#   env             = "prod"                                     # 如果此值不填，则默认读取环境变量中的CTYUN_ENV
# }
#
#
# resource "ctyun_vpc" "vpc_test2" {
#   provider    = ctyun.test_accpet
#   name        = "tf-vpc-for-peer_connect"
#   cidr        = "192.168.0.0/16"
#   description = "terraform测试使用"
#   enable_ipv6 = true
# }

data "ctyun_vpc_route_tables" "route_table_test" {
  vpc_id = ctyun_vpc.vpc_test.id
}