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

resource "ctyun_express_connect" "example" {
  name        = "express_connect_dependence"
  description = "云间高速example专用"
}