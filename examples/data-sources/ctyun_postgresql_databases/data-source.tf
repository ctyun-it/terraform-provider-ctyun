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

# instance id 和 name 根据实际情况进行替换
data "ctyun_postgresql_databases" "exmaples" {
  instance_id = "0e7ed95b886145159c286361463370b5"
}