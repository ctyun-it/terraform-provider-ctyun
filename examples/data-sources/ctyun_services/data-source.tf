terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
}

# 查询服务
data "ctyun_services" "ctyun_services_test" {
  type = "region"
  name = "VPN"
}

output "ctyun_service_test" {
  value = data.ctyun_services.ctyun_services_test.services
}
