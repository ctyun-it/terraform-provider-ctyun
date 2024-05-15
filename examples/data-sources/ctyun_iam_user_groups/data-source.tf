terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
}

data "ctyun_iam_user_groups" "iam_user_groups_test" {
  name      = "terraform"
  page_size = 1000
  page_no   = 1
}

output "iam_user_groups_test" {
  value = data.ctyun_iam_user_groups.iam_user_groups_test.groups
}
