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

resource "ctyun_ec_cloud_gateway" "cloud_gateway_hgh7" {
  ec_id       = ctyun_express_connect.example.id
  name        = "cloud_gateway_hgh7"
  description = "云间高速开发测试专用"
  region_name = "南昌5"
}