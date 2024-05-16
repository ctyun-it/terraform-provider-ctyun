terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

resource "ctyun_enterprise_project" "enterprise_project_test" {
  name        = "my_enterprise_project"
  description = "terraform创建的企业项目"
}