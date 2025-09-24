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

resource "ctyun_eip" "eip_test2" {
 name                = "eip-test2"
 bandwidth           = 10
 cycle_type          = "on_demand"
 demand_billing_type = "bandwidth"
}

resource "ctyun_redis_association_eip" "test" {
  eip_address = ctyun_eip.eip_test2.address
  instance_id = "d59e17a10dda4105936b7e3ede290ba5"
}