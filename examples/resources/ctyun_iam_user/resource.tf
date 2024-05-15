terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
}

resource "ctyun_iam_user" "iam_user_test" {
  email          = "k2mn05@qq.com"
  phone          = "17306692771"
  name           = "Mddi3"
  password       = "P@ssW0rd_!"
  description    = "测试创建账号111"
  user_group_ids = [
    "6edf8a6a9b09442295206feef0d39132"
  ]
}