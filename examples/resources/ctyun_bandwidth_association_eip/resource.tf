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

resource "ctyun_bandwidth_association_eip" "bandwidth_association_eip_test" {
  bandwidth_id = "bandwidth-at2yy664m5"
  eip_id       = "eip-p9qvl63yt6"
}