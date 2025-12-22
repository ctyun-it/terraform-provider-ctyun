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


resource "ctyun_private_zone" "private_zone_example" {
  name          = "test-zone.example.com"
  description   = "terraform dns用例"
  proxy_pattern = "zone"
  ttl           = 300
  vpc_id_list   = [ctyun_vpc.vpc_test.id]
}

resource "ctyun_private_zone_record" "example" {
  zone_id     = ctyun_private_zone.private_zone_example.id
  type        = "A"
  value_list  = ["192.168.1.100"]
  ttl         = 300
  name        = "example-record"
  description = "Example private zone record"
  enabled     = true
}