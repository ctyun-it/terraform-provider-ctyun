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

resource "ctyun_ipv6_gateway" "test" {
  vpc_id = "vpc-d6icy1o968"
}