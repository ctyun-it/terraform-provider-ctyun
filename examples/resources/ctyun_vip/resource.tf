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
resource "ctyun_subnet" "subnet_test" {
  vpc_id      = "vpc-d7zxz8j05c"
  name        = "subnet-test"
  cidr        = "10.0.0.0/8"
  description = "terraform测试使用"
  dns = [
    "114.114.114.114",
    "8.8.8.8"
  ]
  enable_ipv6 = true
}

resource "ctyun_vip" "example" {
  vpc_id     = ctyun_subnet.subnet_test.vpc_id
  subnet_id  = ctyun_subnet.subnet_test.id
  ip_address = "192.168.1.100"
  vip_type   = "v4"
  project_id = "0"
}