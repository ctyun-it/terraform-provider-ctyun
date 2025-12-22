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

resource "ctyun_express_connect" "express_connect_example" {
  name        = "express_connect_example"
  description = "云间高速examples专用"

}
resource "ctyun_ec_cloud_gateway" "cloud_gateway_example" {
  ec_id       = ctyun_express_connect.express_connect_example.id
  name        = "cloud_gateway_example"
  description = "云间高速examples专用"
}
resource "ctyun_vpc" "vpc_test_for_instance" {
  name        = "tf-vpc-for-instance"
  cidr        = "192.168.0.0/16"
  description = "terraform-ec vpc instance测试使用"
  enable_ipv6 = true
  # region_id   = "200000002401"
}


resource "ctyun_subnet" "subnet_test_for_instance" {
  # region_id   = "200000002401"
  vpc_id      = ctyun_vpc.vpc_test_for_instance.id
  name        = "tf-subnet-for-vpc_instance"
  cidr        = "192.168.1.0/24"
  description = "terraform-ec vpc instance测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_ec_vpc_instance" "instance_test" {
  # region_id   = "200000002401"
  ec_id       = ctyun_express_connect.express_connect_example.id
  cgw_id      = ctyun_ec_cloud_gateway.cloud_gateway_example.id
  rtb_id      = ctyun_ec_cloud_gateway.cloud_gateway_example.rtb_id
  vpc_id      = ctyun_vpc.vpc_test_for_instance.id
  route_learn = 1
  route_sync  = 1
  subnets     = [ctyun_subnet.subnet_test_for_instance.id]
}

