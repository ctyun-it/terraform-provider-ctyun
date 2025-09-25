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

resource "ctyun_ebs_backup_policy" "test" {
    name           = "test"
    cycle_type            = "week"
    cycle_week            = "0,2,6"
    time                  = "1,20"
    status                = 1
    retention_type        = "num"
    retention_num         = 20
    full_backup_interval  = -1
    adv_retention_status  = true

    # 当启用高级保留策略且retention_type为num时，可以配置高级保留策略
    adv_retention {
     adv_day = 3
    }
}


