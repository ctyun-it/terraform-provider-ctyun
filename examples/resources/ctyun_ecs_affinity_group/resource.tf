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

resource "ctyun_ecs_affinity_group" "test" {
  affinity_group_name = "tf-test-group"
  affinity_group_policy = "anti-affinity"
}
