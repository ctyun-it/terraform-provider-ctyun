resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-opensearch"
  cidr        = "192.168.0.0/16"
  description = "terraform 测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-opensearch"
  cidr        = "192.168.0.0/24"
  description = "terraform 测试使用"
  dns         = ["8.8.8.8", "8.8.4.4"]
  enable_ipv6 = true
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-opensearch"
  description = "terraform 测试使用"
  lifecycle {
    prevent_destroy = false
  }
}

data "ctyun_zones" "test" {

}

locals {
  az_name      = data.ctyun_zones.test.zones[0]
  zone_list    = toset([local.az_name])
  flavor_name  = "esearch-4c16g"
  storage_type = "SSD-genric"
}
