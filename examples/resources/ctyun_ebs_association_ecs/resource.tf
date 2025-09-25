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
resource "ctyun_ebs_association_ecs" "ebs_association_ecs_test" {
  ebs_id      = "86517323-d3f5-48cb-b278-a124e98fbc3d"
  instance_id = "0b9897c8-ff01-42b4-c4c2-6a427d8b2e9a"
}