terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 创建带宽大小为1Mbps的弹性ip，付费模式为包周期1个月
resource "ctyun_eip" "eip_test1" {
  name        = "eip-test1"
  bandwidth   = 1
  cycle_type  = "month"
  cycle_count = 1
}

## 创建带宽大小为10Mbps的弹性ip，付费模式为按带宽计费
#resource "ctyun_eip" "eip_test2" {
#  name                = "eip-test2"
#  bandwidth           = 10
#  cycle_type          = "on_demand"
#  demand_billing_type = "bandwidth"
#  region_id           = "200000002527"
#  project_id          = "4f5ef15300724760af59b37cf6409f45"
#}
#
## 创建带宽大小为5Mbps的弹性ip，付费模式为按流量计费
#resource "ctyun_eip" "eip_test3" {
#  name                = "eip-test3"
#  bandwidth           = 5
#  cycle_type          = "on_demand"
#  demand_billing_type = "upflowc"
#}