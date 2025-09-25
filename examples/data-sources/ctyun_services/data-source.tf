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

# 查询服务
data "ctyun_services" "ctyun_services_test" {
  type = "region"
  name = "VPN"
}

output "ctyun_service_test" {
  value = data.ctyun_services.ctyun_services_test.services
}
