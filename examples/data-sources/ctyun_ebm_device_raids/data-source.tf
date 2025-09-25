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

data "ctyun_ebm_device_raids" "test" {
  region_id = "200000001852"
  az_name = "cn-huabei2-tj1A-public-ctcloud"
  device_type = "physical.s5.2xlarge4"
  volume_type = "system"
}

output "ctyun_ebm_device_raids_test" {
  value = data.ctyun_ebm_device_raids.test
}

