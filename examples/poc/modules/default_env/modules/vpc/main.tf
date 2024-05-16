terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

resource "ctyun_vpc" "vpc_test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}