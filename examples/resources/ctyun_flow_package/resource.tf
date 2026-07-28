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

# 创建1个月，10G的共享流量包
resource "ctyun_flow_package" "flow_package_test1" {
  cycle_type  = "month"
  spec = 10
}