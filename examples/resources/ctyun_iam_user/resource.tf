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

variable "password" {
  type      = string
  sensitive = true
}

variable "email" {
  type      = string
  sensitive = true
}

variable "phone" {
  type      = string
  sensitive = true
}

resource "ctyun_iam_user" "iam_user_test" {
  email          = var.email
  phone          = var.phone
  name           = "Mddi3"
  password       = var.password
  description    = "测试创建账号111"
  user_group_ids = [
    "6edf8a6a9b09442295206feef0d39132"
  ]
}