terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
}

resource "ctyun_ebs" "ebs_test" {
  name       = "ebs-test"
  mode       = "vbd"
  type       = "sata"
  size       = 60
  cycle_type = "on_demand"
}