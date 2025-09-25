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

resource "ctyun_elb_acl" "test" {
  name = "tf_acl"
  source_ips = ["127.0.0.1/32","192.168.0.0/16","192.168.10.0"]
}
