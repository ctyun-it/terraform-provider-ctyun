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

resource "ctyun_ebs_backup_repo" "test" {
  name =  "tf-test-ctyun_ebs_backup_repo"
  size = 100
  cycle_count = "5"
  cycle_type  = "month"
}
