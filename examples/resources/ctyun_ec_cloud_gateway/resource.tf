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


resource "ctyun_ec_cloud_gateway" "example" {
  ec_id       = ctyun_express_connect.example.id
  name        = "example"
  description = "云间高速example专用"
  region_id   = "200000003329"
  region_name = "cn-zj-hgh7-1a-public-ctcloud"
}
