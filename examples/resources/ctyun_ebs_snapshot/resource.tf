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

resource "ctyun_ebs" "ebs_test" {
  name       = "ebs-test"
  mode       = "vbd"
  type       = "sata"
  size       = 60
  cycle_type = "on_demand"
}

resource "ctyun_ecs_snapshot" "test" {
  name = "tf-test-group"
  disk_id = ctyun_ebs.ebs_test.id
  retention_policy = "forever"
}
