
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


resource "ctyun_mysql_param_template" "examples" {
  name        = "template_example_1"
  engine      = "5.7"
  description = "param template test"
}