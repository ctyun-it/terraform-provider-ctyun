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

resource "ctyun_ebs_backup" "test" {
  repository_id = "0cd13a89-5ada-42a7-95e8-60fb9705eecc"
  disk_id = "f16dfc3f-7375-4831-af16-a4cbd060ec89"
  name  = "test"
  full_backup = false
}


