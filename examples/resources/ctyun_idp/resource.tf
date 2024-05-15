terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
}

resource "ctyun_idp" "idp_test" {
  file        = file("./metadata-saml-idp.xml")
  file_name   = "metadata-saml-idp.xml"
  name        = "minchiang的测试哦"
  type        = "iam"
  protocol    = "saml"
  description = "terraform测试提供商111"
}