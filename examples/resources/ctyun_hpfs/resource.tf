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

resource "ctyun_hpfs" "test" {
  sfs_protocol = "hpfs"
  cycle_type = "on_demand"
  sfs_name = "hpfs-test"
  sfs_size = 512
}

