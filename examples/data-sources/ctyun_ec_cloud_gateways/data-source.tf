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

data "ctyun_ec_cloud_gateways" "example" {
  ec_id = "49410d6d-fd53-48b3-9f78-cb28da38d7be"
}

output "ctyun_ec_cloud_gateways_example" {
  value = data.ctyun_ec_cloud_gateways.example
}