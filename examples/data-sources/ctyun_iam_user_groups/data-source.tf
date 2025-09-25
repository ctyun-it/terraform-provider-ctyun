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

data "ctyun_iam_user_groups" "iam_user_groups_test" {
  name      = "terraform"
  page_size = 1000
  page_no   = 1
}

output "iam_user_groups_test" {
  value = data.ctyun_iam_user_groups.iam_user_groups_test.groups
}
