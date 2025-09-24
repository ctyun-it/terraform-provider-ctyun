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


resource "ctyun_sfs_permission_rule" "sfs_permission_rule_test" {
  permission_group_fuid    =  ctyun_sfs_permission_group.sfs_permission_group_test.id
  auth_addr                = "192.168.1.0/24"
  rw_permission            = "ro"
  permission_rule_priority = 200
}
