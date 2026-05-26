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

data "ctyun_postgresql_param_templates" "param_templates" {

}

resource "ctyun_postgresql_param_template" "example" {
  name               = "pgsql_param_template"
  source_template_id = data.ctyun_postgresql_param_templates.param_templates.parameter_templates[0].id
  description        = "pgsql参数模板"
}


