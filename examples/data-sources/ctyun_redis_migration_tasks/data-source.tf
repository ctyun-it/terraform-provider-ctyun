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

data "ctyun_redis_migration_tasks" "test" {
  id        = "459d960db7c74a3a9ddc8e21cec53597" //任务id  可选
  page_num  = 1
  page_size = 10
}