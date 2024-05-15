terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
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