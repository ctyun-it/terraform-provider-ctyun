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

data "ctyun_private_nat_cidrs" "example" {
  nat_gateway_id = "natgw-asdsmh8scy"
}

output "ctyun_private_nat_cidrs_example" {
  value = data.ctyun_private_nat_cidrs.example
}