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

data "ctyun_hpfs_instances" "test" {
  sfs_status = "available"
  sfs_protocol = "hpfs"
  az_name = "cn-huadong1-jsnj1A-public-ctcloud"
}
