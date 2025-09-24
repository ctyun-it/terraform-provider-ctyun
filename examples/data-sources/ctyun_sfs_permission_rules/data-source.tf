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

resource "ctyun_sfs_permission_group" "sfs_permission_group_test" {
  name = "permission-group_example"
  description = "创建sfs规则组"
}

data "ctyun_sfs_permission_rules" "test" {
  permission_group_fuid = ctyun_sfs_permission_group.sfs_permission_group_test.id
}
