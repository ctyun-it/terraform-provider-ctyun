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

data "ctyun_ecs_flavors" "ecs_flavor_test" {
  cpu    = 4
  ram    = 8
  arch   = "x86"
  series = "C"
  type   = "CPU_C7"
}

resource "ctyun_ccse_node_pool" "example" {
  cluster_id               = "dd92f3a6b034431bb7dceb849aed1220"
  name           = "default-pool"
  cycle_type              = "month"
  cycle_count = 1
  auto_renew = true
  instance_type            = "ecs"
  mirror_id                = "3f80d8c0-8eb5-4afa-a506-13ba68b61872"
  mirror_type              = 1
  key_pair_name           = "KeyPair-de15"
  use_affinity_group       = true
  affinity_group_id      = "f8b18511-4327-4c3f-9373-c6d661889fcb"
  item_def_name            = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].name
  max_pod_num              = 110
  az_infos = [
    {
      az_name = "cn-huadong1-jsnj1A-public-ctcloud"
    }
  ]
  sys_disk = {
    type = "SATA"
    size = 300
  }

  data_disks = [
    {
      type = "SSD"
      size = 4000
    }
  ]
}

# 裸金属节点池（physical.s5.2xlarge4不支持云硬盘）
# resource "ctyun_ccse_node_pool" "example" {
#   cluster_id               = "dd92f3a6b034431bb7dceb849aed1220"
#   name           = "default-pool1"
#   cycle_type              = "month"
#   cycle_count = 1
#   auto_renew = true
#   instance_type            = "ebm"
#   mirror_name             = "CTyunOS23.01@cpu_ccse_img_4.0_09"
#   mirror_type              = 1
#   key_pair_name           = "KeyPair-de15"
#   use_affinity_group       = true
#   affinity_group_id      = "f8b18511-4327-4c3f-9373-c6d661889fcb"
#   item_def_name            = "physical.s5.2xlarge4"
#   max_pod_num              = 110
#   az_infos = [
#     {
#       az_name = "cn-huadong1-jsnj1A-public-ctcloud"
#     }
#   ]
# }