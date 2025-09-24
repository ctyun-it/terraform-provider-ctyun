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

data "ctyun_ccse_template_market" "test" {
  tpl_name = "elasticsearch"
  tpl_version = "7.10.2"
  values_type = "YAML"
}

resource "ctyun_ccse_template_instance" "test_template" {
  cluster_id = "f440fe8c26c94dd88adea41a08a5353d"
  tpl_name = "elasticsearch"
  tpl_version = "7.10.2"
  name = "tf-ccse-12"
  namespace = "default"
  values_yaml = data.ctyun_ccse_template_market.test.values
}