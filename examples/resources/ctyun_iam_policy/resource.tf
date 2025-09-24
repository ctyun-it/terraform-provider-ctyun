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

resource "ctyun_iam_policy" "iam_policy_test" {
  name        = "terraform_policy_test1"
  description = "terraform测试新建策略"
  range       = "region"
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
      },
      {
        effect = "allow"
        action = ["kms:cmk:list"]
      }
    ]
  }
}
