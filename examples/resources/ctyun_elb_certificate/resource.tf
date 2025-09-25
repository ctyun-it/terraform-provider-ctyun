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

resource "ctyun_elb_certificate" "certificate_test" {
  name        = "tf_elb_certifiate"
  type        = "Server"
  certificate = "xxxxxxxx"
  private_key = "xxxxxxxx"
}
