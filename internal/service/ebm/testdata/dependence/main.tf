resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-ebm"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-ebm"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
  enable_ipv6 = true
  type = "common"
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-ebm"
  description = "terraform测试使用"
}

resource "ctyun_security_group" "security_group_test2" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-ebm2"
  description = "terraform测试使用"
}

data "ctyun_zones" "test" {

}

locals {
  device_type1 = "physical.s5.2xlarge4"      // az2、有本地盘、弹性、不支持云硬盘
  device_type2 = "physical.s5.2xlarge1"      // az2、无本地盘、弹性、支持云硬盘
  az2 = data.ctyun_zones.test.zones[1]
}

data "ctyun_ebm_device_raids" "system_raid" {
  az_name = local.az2
  device_type = local.device_type1
  volume_type = "system"
}

data "ctyun_ebm_device_raids" "data_raid" {
  az_name = local.az2
  device_type = local.device_type1
  volume_type = "data"
}

data "ctyun_ebm_device_images" "test" {
  az_name = local.az2
  device_type = local.device_type1
  os_type = "linux"
  image_type = "standard"
}

locals {
  system_raid_id = data.ctyun_ebm_device_raids.system_raid.raids[0].uuid
  data_raid_id = data.ctyun_ebm_device_raids.data_raid.raids[0].uuid
}


data "ctyun_ebm_device_images" "dependence" {
  device_type = local.device_type2
  az_name = local.az2
  os_type = "linux"
  image_type = "standard"
}

resource "ctyun_ebs" "ebs_test" {
  az_name   = local.az2
  name       = "tf-ebs-for-ebm"
  mode       = "vbd"
  type       = "sata"
  size       = 60
  cycle_type = "on_demand"
}

resource "ctyun_eip" "eip_test" {
  name                = "tf-eip-for-ebm"
  bandwidth           = 1
  cycle_type          = "on_demand"
  demand_billing_type = "upflowc"
}

resource "ctyun_ebm" "ebm_test" {
  az_name   = local.az2
  instance_name = "tf-ebm-for-ebm"
  hostname = "tf-ebm-for-ebm"
  password = var.password
  eip_id = ctyun_eip.eip_test.id
  cycle_type = "on_demand"
  device_type = local.device_type2
  image_uuid = data.ctyun_ebm_device_images.dependence.images[0].image_uuid
  security_group_ids = [ctyun_security_group.security_group_test.id]
  vpc_id = ctyun_vpc.vpc_test.id
  system_disk_size = 100
  system_disk_type = "sata"
  subnet_id = ctyun_subnet.subnet_test.id
}

variable "password" {
  type      = string
  sensitive = true
}
