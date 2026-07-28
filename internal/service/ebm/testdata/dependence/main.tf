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
    "8.8.8.8"
  ]
  enable_ipv6 = true
  type = "common"
}

resource "ctyun_subnet" "subnet_ebm" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-standard"
  cidr        = "192.168.2.0/24"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8"
  ]
  enable_ipv6 = true
  type = "cbm"
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

data "ctyun_ebm_device_types" "test" {
  for_each = { for az in data.ctyun_zones.test.zones: az => az }
  az_name = each.value
}

locals {
  # 筛选全部有余量的规格
  all_device_types = flatten([
    for device_type_inst in values(data.ctyun_ebm_device_types.test) : [
      for dt in device_type_inst.device_types : dt
      if dt.available == true # 核心过滤条件
    ]
  ])

  # 支持云盘的弹性裸金属
  cloud_boot_true_list = [for dt in local.all_device_types : dt if dt.cloud_boot == true && dt.smart_nic_exist == true]
  first_cloud_boot_true = length(local.cloud_boot_true_list) > 0 ? local.cloud_boot_true_list[0] : null

  # 不支持云盘的弹性裸金属
  cloud_boot_false_list = [for dt in local.all_device_types : dt if dt.cloud_boot == false && dt.smart_nic_exist == true]
  first_cloud_boot_false = length(local.cloud_boot_false_list) > 0 ? local.cloud_boot_false_list[0] : null

  # 标准裸金属
  standard_list = [for dt in local.all_device_types : dt if dt.smart_nic_exist == false]
  standard = length(local.standard_list) > 0 ? local.standard_list[0] : null
}

locals {
  device_type1 = local.first_cloud_boot_false != null ? local.first_cloud_boot_false.device_type : ""
  device_type2 = local.first_cloud_boot_true != null ? local.first_cloud_boot_true.device_type : ""
  standard_device_type = local.standard != null ? local.standard.device_type : ""

  first_cloud_boot_true_az = local.first_cloud_boot_true != null ? local.first_cloud_boot_true.az_name : ""
  first_cloud_boot_false_az = local.first_cloud_boot_false != null ? local.first_cloud_boot_false.az_name : ""
  standard_az = local.standard != null ? local.standard.az_name : ""
}

data "ctyun_ebm_device_raids" "system_raid" {
  az_name = local.first_cloud_boot_false_az
  device_type = local.device_type1
  volume_type = "system"
}

data "ctyun_ebm_device_raids" "data_raid" {
  az_name = local.first_cloud_boot_false_az
  device_type = local.device_type1
  volume_type = "data"
}

data "ctyun_ebm_device_raids" "system_raid_standard" {
  az_name = local.standard_az
  device_type = local.standard_device_type
  volume_type = "system"
}

data "ctyun_ebm_device_raids" "data_raid_standard" {
  az_name = local.standard_az
  device_type = local.standard_device_type
  volume_type = "data"
}

data "ctyun_ebm_device_images" "test" {
  az_name = local.first_cloud_boot_false_az
  device_type = local.device_type1
  os_type = "linux"
  image_type = "standard"
}

data "ctyun_ebm_device_images" "standard_test" {
  az_name = local.standard_az
  device_type = local.standard_device_type
  os_type = "linux"
  image_type = "standard"
}


locals {
  system_raid_id = data.ctyun_ebm_device_raids.system_raid.raids[0].uuid
  data_raid_id = data.ctyun_ebm_device_raids.data_raid.raids[0].uuid

  standard_system_raid_id = data.ctyun_ebm_device_raids.system_raid_standard.raids[0].uuid
  standard_data_raid_id = data.ctyun_ebm_device_raids.data_raid_standard.raids[0].uuid
}


data "ctyun_ebm_device_images" "dependence" {
  device_type = local.device_type2
  az_name = local.first_cloud_boot_true_az
  os_type = "linux"
  image_type = "standard"
}

resource "ctyun_ebs" "ebs_test" {
  az_name = local.first_cloud_boot_true_az
  name       = "tf-ebs-for-ebm"
  mode       = "vbd"
  type       = "SATA"
  size       = 60
  cycle_type = "on_demand"
}

resource "ctyun_ebm" "ebm_test" {
  az_name = local.first_cloud_boot_true_az
  instance_name = "tf-ebm-for-ebm"
  hostname = "tf-ebm-for-ebm"
  password = var.password
  bandwidth = 2
  cycle_type = "on_demand"
  device_type = local.device_type2
  image_uuid = data.ctyun_ebm_device_images.dependence.images[0].image_uuid
  security_group_ids = [ctyun_security_group.security_group_test.id]
  vpc_id = ctyun_vpc.vpc_test.id
  system_disk_size = 100
  system_disk_type = "SATA"
  subnet_id = ctyun_subnet.subnet_test.id
}

variable "password" {
  type      = string
  sensitive = true
}