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

resource "ctyun_ipv6_bandwidth_association" "test" {
  ipv6_bandwidth_id = "ipv6-bandwidth-3"
  ipv6 = 3
}