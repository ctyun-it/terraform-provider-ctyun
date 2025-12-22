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

data "ctyun_iam_policies" "policies" {
  page_no   = 1
  page_size = 100
}


data "ctyun_iam_policies" "policy" {
  policy_id = "c6b9d3dc3d3c426c833469760ee15e61"
}