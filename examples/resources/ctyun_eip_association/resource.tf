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

resource "ctyun_eip_association" "eip_association_test2" {
  eip_id      = "eip-nl78g1t31o"
  instance_id = "fd94fbe2-26b2-5dbb-5deb-65b4167ca28e"
  region_id   = "200000002527"
  project_id  = "4f5ef15300724760af59b37cf6409f45"
}