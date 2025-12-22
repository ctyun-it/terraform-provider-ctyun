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

data "ctyun_net_tagss" "example" {
  resource_type = "vpc"
  resource_id   = "vpc-asdsmh8scy"
}

output "ctyun_net_tagss_example" {
  value = data.ctyun_net_tagss.example
}