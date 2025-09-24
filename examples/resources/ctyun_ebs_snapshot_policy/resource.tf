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


resource "ctyun_ebs_snapshot_policy" "test" {
    name           = "test"
    repeat_weekdays            = "0,1,2"
    repeat_times            = "0,1,2"
    retention_time        = 2
    is_enabled  = true
}

