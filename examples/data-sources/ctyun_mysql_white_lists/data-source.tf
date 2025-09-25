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


data "ctyun_mysql_white_lists" "test" {
  prod_inst_id = "xxxxxxxxxxxxxx"   # 数据库实例ID
}

