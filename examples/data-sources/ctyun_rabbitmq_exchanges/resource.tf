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

# 全部
data "ctyun_rabbitmq_exchanges" "all" {
  instance_id = "8ccc8af2e6704080a72548735a081660"
}

# 指定vhost
data "ctyun_rabbitmq_exchanges" "test" {
  instance_id = "8ccc8af2e6704080a72548735a081660"
  vhost       = "/"
}