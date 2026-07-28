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

resource "ctyun_iam_user_ak" "ak_test" {
  user_id = "c2c8827c8ca4433a97e8e133a41c33b9"
  enabled = false
}