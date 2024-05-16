terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id = var.vpc_id
  name   = var.subnet_name
  cidr   = var.subnet_cidr
  dns    = var.subnet_dns
}