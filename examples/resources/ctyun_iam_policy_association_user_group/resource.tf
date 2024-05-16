terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

resource "ctyun_iam_policy" "iam_policy_global_test" {
  name        = "terraform_policy_global2"
  description = "terraform测试全局策略2"
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
  name        = "terraform_policy_region2"
  description = "terraform测试资源池策略2"
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

# 绑定资全局型的策略
resource "ctyun_iam_policy_association_user_group" "iam_policy_association_user_group_test" {
  user_group_id = "6edf8a6a9b09442295206feef0d39132"
  policy_id     = ctyun_iam_policy.iam_policy_global_test.id
}

# 绑定资源池型的策略
resource "ctyun_iam_policy_association_user_group" "iam_policy_association_user_region_test" {
  user_group_id = "6edf8a6a9b09442295206feef0d39132"
  policy_id     = ctyun_iam_policy.iam_policy_region_test.id
  region_id     = "bb9fdb42056f11eda1610242ac110002"
}