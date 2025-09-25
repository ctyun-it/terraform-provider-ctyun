terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 可参考index.md，在环境变量中配置ak、sk、资源池ID、可用区名称
provider "ctyun" {
  env = "prod"
}


resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-sfs"
  cidr        = "192.168.0.0/16"
  description = "terraform-sfs测试使用"
  enable_ipv6 = true
}

resource "ctyun_sfs_permission_group" "sfs_permission_group_test" {
  name = "permission-group_example"
  description = "创建sfs规则组"
}


resource "ctyun_vpc" "vpc_test1" {
  name        = "tf-vpc-for-sfs1"
  cidr        = "192.168.0.0/16"
  description = "terraform-sfs测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-sfs"
  cidr        = "192.168.1.0/24"
  description = "terraform-sfs测试使用"
  dns = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
}
resource "ctyun_sfs" "sfs_test" {
  sfs_type     = "performance"
  sfs_protocol = "nfs"
  name         = "sfs-example"
  sfs_size     = 500
  cycle_type   = "on_demand"
  az_name      = "cn-huadong1-jsnj1A-public-ctcloud"
  vpc_id       = ctyun_vpc.vpc_test.id
  subnet_id    = ctyun_subnet.subnet_test.id
}


resource "ctyun_sfs_permission_group_association" "sfs_permission_group_association_test" {
  permission_group_fuid = ctyun_sfs_permission_group.sfs_permission_group_test.id
  sfs_uid               = ctyun_sfs.sfs_test.id
  vpc_id                = ctyun_vpc.vpc_test1.id
}
