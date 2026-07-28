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

data "ctyun_prefix_lists" "example" {
  id            = "pl-1234567890abcdef1234567890abcdef"
  query_content = "test-prefix-ipv6-"
}

output "ctyun_prefix_lists_example" {
  value = data.ctyun_prefix_lists.example
}