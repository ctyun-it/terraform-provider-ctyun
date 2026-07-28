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


resource "ctyun_oceanfs_permission_group" "example" {
  name        = "oceanf_pg_example"
  description = "terraform样例"
}



resource "ctyun_oceanfs_permission_rule" "example" {
  permission_group_id = ctyun_oceanfs_permission_group.example.id
  auth_addr           = "192.168.1.0/24"
  rw_permission       = "ro"
  priority            = 1
}

