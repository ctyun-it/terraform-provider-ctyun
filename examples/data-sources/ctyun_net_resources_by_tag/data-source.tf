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

data "ctyun_net_resources_by_tag" "example" {
  label_id = "1c24ddb1ff534de9a4bcd13c3b680b59"
}

output "ctyun_net_resources_by_tag_example" {
  value = data.ctyun_net_resources_by_tag.example
}