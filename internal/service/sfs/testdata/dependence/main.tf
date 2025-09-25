resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-sfs"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_vpc" "vpc_test1" {
  name        = "tf-vpc-for-sfs1"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}


resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-sfs"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_subnet" "subnet_test1" {
  vpc_id      = ctyun_vpc.vpc_test1.id
  name        = "tf-subnet-for-sfs1"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_sfs" "sfs_test" {
  sfs_type     = "capacity"
  sfs_protocol = "nfs"
  name         = "sfs-for-group"
  sfs_size     = 500
  cycle_type   = "on_demand"
  vpc_id       = ctyun_vpc.vpc_test.id
  subnet_id    = ctyun_subnet.subnet_test.id
}

resource "ctyun_sfs_permission_group" "group_test" {
  name = "sfs-test1"
  description = "单元测试1"
}

resource "ctyun_sfs_permission_group" "group_test1" {
  name = "sfs-test2"
  description = "单元测试2"
}