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

# 查找1c1g x86架构的通用型的规格
data "ctyun_ecs_flavors" "ctyun_ecs_flavors_test1" {
  cpu    = 1
  ram    = 1
  arch   = "x86"
  series = "S"
  type   = "CPU_S7"
}

# 通过规格名称查找
data "ctyun_ecs_flavors" "ctyun_ecs_flavors_test2" {
  name = "pi7.4xlarge.4"
}

output "ctyun_ecs_flavor_id1" {
  value = data.ctyun_ecs_flavors.ctyun_ecs_flavors_test1.flavors
}

output "ctyun_ecs_flavor_id2" {
  value = data.ctyun_ecs_flavors.ctyun_ecs_flavors_test2.flavors
}