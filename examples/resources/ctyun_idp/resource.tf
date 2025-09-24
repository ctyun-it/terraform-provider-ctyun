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

resource "ctyun_idp" "idp_test" {
  file        = file("./metadata-saml-idp.xml")
  file_name   = "metadata-saml-idp.xml"
  name        = "minchiang的测试哦"
  type        = "iam"
  protocol    = "saml"
  description = "terraform测试提供商111"
}