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

resource "ctyun_acl" "example" {
  vpc_id             = "vpc-idexample1"
  name               = "example-acl"
  description        = "Example ACL created for demonstration"
  enabled            = true
  apply_to_public_lb = false
}
resource "ctyun_subnet_association_acl" "example" {
  acl_id    = ctyun_acl.example.id
  subnet_id = ctyun_subnet.subnet_test.id
}