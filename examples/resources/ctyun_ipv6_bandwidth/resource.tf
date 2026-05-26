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

resource "ctyun_ipv6_bandwidth" "test" {
  name = "ipv6-bandwidth-3"
  bandwidth = 3
  cycle_type = "month"
  cycle_count = 1
}