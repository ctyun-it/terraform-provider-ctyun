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

data "ctyun_mongodb_white_lists" "example" {
  instance_id = "02d8a872-7d1a-45ab-9bd8-9b158376ba3a"
}

output "ctyun_mongodb_white_lists_example" {
  value = data.ctyun_mongodb_white_lists.example
}