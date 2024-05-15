terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
}

resource "ctyun_iam_policy" "iam_policy_global_test" {
  name        = "terraform_policy_global"
  description = "terraform测试全局策略"
  range       = "global"
  content     = {
    version   = "1.1"
    statement = [
      {
        effect   = "allow"
        action   = ["vpc:vpcs:list", "vpc:vpcs:get"]
        resource = ["*"]
      },
      {
        effect   = "allow"
        action   = ["vpc:vpcs:delete"]
        resource = ["*"]
      }
    ]
  }
}

resource "ctyun_iam_policy" "iam_policy_region_test" {
  name        = "terraform_policy_region"
  description = "terraform测试资源池策略"
  range       = "region"
  content     = {
    version   = "1.1"
    statement = [
      {
        effect   = "allow"
        action   = ["kms:cmk:list"]
        resource = ["*"]
      }
    ]
  }
}

# 绑定全局型的策略
resource "ctyun_iam_policy_association_user" "iam_policy_association_user_global_test" {
  user_id   = "f9c327e495de4b3db569cf3b604e4d76"
  policy_id = ctyun_iam_policy.iam_policy_global_test.id
}

# 绑定资源池型的策略
resource "ctyun_iam_policy_association_user" "iam_policy_association_user_region_test" {
  user_id   = "f9c327e495de4b3db569cf3b604e4d76"
  policy_id = ctyun_iam_policy.iam_policy_region_test.id
  region_id = "bb9fdb42056f11eda1610242ac110002"
}