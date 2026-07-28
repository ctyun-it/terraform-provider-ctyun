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


resource "ctyun_redis_param_template" "test" {
  name         = "testname"
  description  = "Initial Redis template"
  cache_mode   = "ORIGINAL_67"
  sys_template = false

  params = [{
    param_name    = "maxmemory-policy"
    current_value = "allkeys-lru"
  }]
}
