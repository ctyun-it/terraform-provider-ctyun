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


# 跨帐号
resource "ctyun_vpc_peer_connection" "cross_example" {
  # provider = ctyun.default
  project_id     = "0"
  name           = "vpc-peer-conn-example"
  description    = "对等连接样例"
  request_vpc_id = ctyun_vpc.vpc_test.id
  accept_vpc_id  = ctyun_vpc.vpc_test1.id
  accept_email   = "xxxx@chinatelecom.cn"
}


