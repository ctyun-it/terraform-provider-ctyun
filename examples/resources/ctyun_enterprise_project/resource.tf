terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
}

resource "ctyun_enterprise_project" "enterprise_project_test" {
  name        = "my_enterprise_project"
  description = "terraform创建的企业项目"
}