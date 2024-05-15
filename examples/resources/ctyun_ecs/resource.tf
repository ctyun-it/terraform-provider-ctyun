terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
}

provider "ctyun" {
  region_id  = "bb9fdb42056f11eda1610242ac110002"
  project_id = "17a308cdf06a4a7ebfb27070a7b07e18"
  az_name    = "cn-huadong1-jsnj1A-public-ctcloud"
}

data "ctyun_images" "image_test1" {
  name       = "Ubuntu 22.04"
  visibility = "public"
  page_no    = 1
  page_size  = 50
}

data "ctyun_ecs_flavors" "ecs_flavor_test1" {
  cpu    = 1
  ram    = 1
  arch   = "x86"
  series = "S"
  type   = "CPU_S7"
}

# 创建1c1g x86架构的通用型云主机，系统盘使用SATA，40g，计费方式为包周期1个月，到期自动续费
resource "ctyun_ecs" "ecs_test1" {
  name               = "ecs-test"
  flavor_id          = data.ctyun_ecs_flavors.ecs_flavor_test1.flavors[0].id
  image_id           = data.ctyun_images.image_test1.images[0].id
  system_disk_type   = "sata"
  system_disk_size   = 40
  vpc_id             = "vpc-r7kv00qbz5"
  password           = "P@ssW0rd_1"
  cycle_type         = "month"
  cycle_count        = 1
  auto_renew         = true
  subnet_id          = "subnet-f3ktwpsf07"
  security_group_ids = [
  ]
}

#########################################################################################

data "ctyun_images" "image_test2" {
  name       = "Ubuntu 22.04"
  visibility = "public"
  page_no    = 1
  page_size  = 50
  region_id  = "200000002527"
  az_name    = "cn-jx-nc5-jxnc1A-public-ctcloud"
}

data "ctyun_ecs_flavors" "ecs_flavor_test2" {
  cpu       = 2
  ram       = 4
  arch      = "x86"
  series    = "S"
  type      = "CPU_S7"
  region_id = "200000002527"
  az_name   = "cn-jx-nc5-jxnc1A-public-ctcloud"
}

# 创建2c4g x86架构的通用型云主机，系统盘使用SATA，50g，计费方式为按需
resource "ctyun_ecs" "ecs_test2" {
  name               = "ecs-test2"
  flavor_id          = data.ctyun_ecs_flavors.ecs_flavor_test2.flavors[0].id
  image_id           = data.ctyun_images.image_test2.images[0].id
  system_disk_type   = "sata"
  system_disk_size   = 50
  vpc_id             = "vpc-d7zxz8j05c"
  password           = "P@ssW0rd_1"
  cycle_type         = "on_demand"
  subnet_id          = "subnet-5jtwyd0m15"
  security_group_ids = [
    "sg-p5gue21vy8"
  ]
  key_pair_name = "keypair-test"
  region_id     = "200000002527"
  az_name       = "cn-jx-nc5-jxnc1A-public-ctcloud"
  project_id    = "4f5ef15300724760af59b37cf6409f45"
}