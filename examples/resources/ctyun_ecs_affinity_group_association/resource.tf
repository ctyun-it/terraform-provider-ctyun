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

resource "ctyun_ecs_affinity_group_association" "test" {
  instance_id = "ae432721-61bf-45b7-b207-7e3256c1c2d6"
  affinity_group_id = "e9d3239a-207a-4006-aa84-3945265bac27"
}