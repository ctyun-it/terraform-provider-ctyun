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

resource "ctyun_iam_user_group" "user_group_test" {
  name        = "terraform_user_group"
  description = "terraform_user_group用户组"
}