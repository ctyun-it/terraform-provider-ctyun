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

resource "ctyun_mysql_white_list" "test" {
  prod_inst_id = "e5ad1c553e394bc891c5bf8fc58be191"
  group_name = "eip_white"
  group_white_list = ["192.168.1.1", "30.8.7.*"]
}