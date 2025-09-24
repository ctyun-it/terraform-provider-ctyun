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

resource "ctyun_ecs_backup_policy_bind_instances" "test" {
  id = "a4b793881bbd42edaa6a0002900e5819"
  instance_id_list = "ae432721-61bf-45b7-b207-7e3256c1c2d6"
}
