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

resource "ctyun_ecs" "jdxutuzpfr" {
  instance_name      = "tf-test-ecs"
  display_name       = "tf-test-init-ecs"
  flavor_id          = "9b4b5e39-db25-f2c8-3914-76881ee77d5c"
  image_id           = "fa3f3784-34f9-4f6b-80a1-dd173d486bd6"
  system_disk_type   = "sata"
  system_disk_size   = 60
  vpc_id             = "vpc-0ein2p8bs8"
  subnet_id          = "subnet-0oiyrpu8nk"
  key_pair_name      = "tf-keypair-for-ecs"
  cycle_type         = "on_demand"
}

# variable "password" {
#   type      = string
#   sensitive = true
# }
# data "ctyun_images" "image_test1" {
#   name       = "Ubuntu 22.04"
#   visibility = "public"
#   page_no    = 1
#   page_size  = 50
# }
#
# data "ctyun_ecs_flavors" "ecs_flavor_test1" {
#   cpu    = 1
#   ram    = 1
#   arch   = "x86"
#   series = "S"
#   type   = "CPU_S7"
# }
#
# # 创建1c1g x86架构的通用型云主机，系统盘使用SATA，40g，计费方式为包周期1个月，到期自动续费
# resource "ctyun_ecs" "ecs_test1" {
#   instance_name       = "ecs-test"
#   display_name        = "ecs-test"
#   flavor_id           = data.ctyun_ecs_flavors.ecs_flavor_test1.flavors[0].id
#   image_id            = data.ctyun_images.image_test1.images[0].id
#   system_disk_type    = "sata"
#   system_disk_size    = 40
#   vpc_id              = "vpc-r7kv00qbz5"
#   password            = var.password
#   cycle_type          = "month"
#   cycle_count         = 1
#   auto_renew          = true
#   subnet_id           = "subnet-f3ktwpsf07"
#   is_destroy_instance = false
# }
#
# #########################################################################################
#
# data "ctyun_images" "image_test2" {
#   name       = "Ubuntu 22.04"
#   visibility = "public"
#   page_no    = 1
#   page_size  = 50
#   region_id  = "200000002527"
#   az_name    = "cn-jx-nc5-jxnc1A-public-ctcloud"
# }
#
# data "ctyun_ecs_flavors" "ecs_flavor_test2" {
#   cpu       = 2
#   ram       = 4
#   arch      = "x86"
#   series    = "S"
#   type      = "CPU_S7"
#   region_id = "200000002527"
#   az_name   = "cn-jx-nc5-jxnc1A-public-ctcloud"
# }
#
# # 创建2c4g x86架构的通用型云主机，系统盘使用SATA，50g，计费方式为按需
# resource "ctyun_ecs" "ecs_test2" {
#   instance_name      = "ecs-test2"
#   display_name       = "ecs-test2"
#   flavor_id          = data.ctyun_ecs_flavors.ecs_flavor_test2.flavors[0].id
#   image_id           = data.ctyun_images.image_test2.images[0].id
#   system_disk_type   = "sata"
#   system_disk_size   = 50
#   vpc_id             = "vpc-d7zxz8j05c"
#   password           = var.password
#   cycle_type         = "on_demand"
#   subnet_id          = "subnet-5jtwyd0m15"
#   security_group_ids = [
#     "sg-p5gue21vy8"
#   ]
#   key_pair_name = "keypair-test"
#   region_id     = "200000002527"
#   az_name       = "cn-jx-nc5-jxnc1A-public-ctcloud"
#   project_id    = "4f5ef15300724760af59b37cf6409f45"
# }