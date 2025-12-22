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
  description = "云间高速开发测试专用"

}

# 创建云间高速带宽包
resource "ctyun_ec_packet" "example" {
  ec_id       = ctyun_express_connect.example.id
  name        = "example-ec-packet"
  bandwidth   = 10
  cycle_type  = "month"
  cycle_count = 1
  area_a      = "china"
  area_b      = "china"
}

