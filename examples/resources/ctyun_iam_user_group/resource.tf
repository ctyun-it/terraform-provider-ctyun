terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
}

resource "ctyun_iam_user_group" "user_group_test" {
  name        = "terraform_user_group"
  description = "terraform_user_group用户组"
}