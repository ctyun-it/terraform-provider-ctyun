
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

resource "ctyun_express_connect" "express_connect_dependence" {
  name        = "express_connect_dependence"
  description = "云间高速example专用"

}
resource "ctyun_ec_cloud_gateway" "cloud_gateway_dependence" {
  ec_id       = ctyun_express_connect.express_connect_dependence.id
  name        = "cloud_gateway_dependence"
  description = "云间高速example专用"
}

resource "ctyun_ec_cloud_gateway" "cloud_gateway_huhehaote3" {
  ec_id       = ctyun_express_connect.express_connect_dependence.id
  name        = "cloud_gateway_xinan1"
  description = "云间高速example专用"
  region_id   = "200000003573"
  region_name = "cn-nm-het3-1a-public-ctcloud"
}

resource "ctyun_ec_cloud_gateway" "cloud_gateway_wulumuqi7" {
  ec_id       = ctyun_express_connect.express_connect_dependence.id
  name        = "cloud_gateway_hgh7"
  description = "云间高速example专用"
  region_id   = "200000004098"
  region_name = "cn-xj-urc7-1a-public-ctcloud"
}

resource "ctyun_ec_packet" "packet_test" {
  ec_id       = ctyun_express_connect.express_connect_dependence.id
  name        = "packet_region_peer_test"
  bandwidth   = 10
  cycle_type  = "month"
  cycle_count = 1
}

resource "ctyun_ec_region_peer" "region_peer_test" {
  name        = "region_peer_test"
  ec_id       = ctyun_express_connect.express_connect_dependence.id
  src_cgw_id  = ctyun_ec_cloud_gateway.cloud_gateway_huhehaote3.id
  dst_cgw_id  = ctyun_ec_cloud_gateway.cloud_gateway_wulumuqi7.id
  packet_id   = ctyun_ec_packet.packet_test.id
  rate        = 1
  route_learn = 1
}