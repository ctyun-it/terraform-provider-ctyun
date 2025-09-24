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

# 创建一个按需，大小为5Mbit/s的带宽
resource "ctyun_bandwidth" "bandwidth_test1" {
  name       = "bandwidth-test1"
  cycle_type = "on_demand"
  bandwidth  = 5
}

# 创建一个包年，大小为10Mbit/s的带宽
resource "ctyun_bandwidth" "bandwidth_test2" {
  name        = "bandwidth-test2"
  cycle_type  = "year"
  bandwidth   = 10
  cycle_count = 1
}