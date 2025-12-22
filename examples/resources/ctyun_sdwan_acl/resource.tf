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

resource "ctyun_sdwan_acl" "acl_test" {
  name = "acl_test"

  rules = [{
    direction      = "in"
    protocol       = "udp"
    ip_version     = "IPv4"
    dst_cidr       = "10.0.0.0/16"
    dst_port_range = "-1/-1"
    priority       = 100
    action         = "allow"
    src_cidr       = "10.0.0.0/16"
    src_port_range = "-1/-1"
  }]
}

