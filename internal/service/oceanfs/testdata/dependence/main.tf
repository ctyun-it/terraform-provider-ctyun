resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-oceanfs-1${local.random_string}"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_vpc" "vpc_test1" {
  name        = "tf-vpc-for-oceanfs-2${local.random_string}"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}


resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-oceanfs-1${local.random_string}"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_subnet" "subnet_test1" {
  vpc_id      = ctyun_vpc.vpc_test1.id
  name        = "tf-subnet-for-oceanfs-2${local.random_string}"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}


resource "ctyun_oceanfs_permission_group" "test" {
  name        = "pg-for-tf"
  description = "terraform测试使用"
}

resource "ctyun_oceanfs_permission_group" "test1" {
  name        = "pg-for-tf1"
  description = "terraform测试使用1"
}

resource "ctyun_oceanfs" "test" {
  protocol      = "nfs"
  name         = "oceanfs-for-tf"
  size          = 100
  cycle_type   = "on_demand"
  vpc_id       = ctyun_vpc.vpc_test.id
  subnet_id    = ctyun_subnet.subnet_test.id
}

locals {
  random_string = substr(replace(lower(sha256(timestamp())), "/[^a-z0-9]/", ""), 0, 5)
}