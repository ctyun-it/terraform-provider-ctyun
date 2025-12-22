data "ctyun_vpcs" "vpc_test" {
  page_size = 50
}

locals {
  vpcs = [for vpc in data.ctyun_vpcs.vpc_test.vpcs : vpc if vpc.name == "tf-vpc-for-paas"]
  data_vpc_id = length(local.vpcs) > 0 ? local.vpcs[0].vpc_id : ""
}

resource "ctyun_vpc" "vpc_test" {
  count    = local.data_vpc_id == "" ? 1 : 0
  name        = "tf-vpc-for-paas"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

locals {
  real_vpc_id = local.data_vpc_id == "" ? try(ctyun_vpc.vpc_test[0].id, "") : local.data_vpc_id
}

data "ctyun_subnets" "subnet_test" {
  vpc_id = local.real_vpc_id
}

locals {
  subnets = [for subnet in data.ctyun_subnets.subnet_test.subnets : subnet if subnet.name == "tf-subnet-for-paas"]
  data_subnet_id = length(local.subnets) > 0 ? local.subnets[0].subnet_id : ""
}

resource "ctyun_subnet" "subnet_test" {
  count    = local.data_vpc_id == "" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-subnet-for-paas"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns         = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

locals {
  real_subnet_id = local.data_subnet_id == "" ? try(ctyun_subnet.subnet_test[0].id, "") : local.data_subnet_id
}

data "ctyun_security_groups" "security_group_test" {
  vpc_id = local.real_vpc_id
}

locals {
  security_groups = [for security_group in data.ctyun_security_groups.security_group_test.security_groups : security_group if security_group.name == "tf-sg-for-paas"]
  data_security_group_id = length(local.security_groups) > 0 ? local.security_groups[0].security_group_id : ""
}

resource "ctyun_security_group" "security_group_test" {
  count    = local.data_vpc_id == "" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-sg-for-paas"
  description = "terraform测试使用"
  lifecycle {
    prevent_destroy = true
  }
}

locals {
  real_security_group_id = local.data_security_group_id == "" ? try(ctyun_security_group.security_group_test[0].id, "") : local.data_security_group_id
}


data "ctyun_ecs_instances" "ecs_test" {
  page_size = 50
}

locals {
  instances = [for instance in data.ctyun_ecs_instances.ecs_test.instances : instance if instance.instance_name == "tf-test-ecs-for-paas"]
  ecs_instance_id = length(local.instances) > 0 ? local.instances[0].id : ""
}

data "ctyun_images" "image_test" {
  name       = "CentOS Linux 8.4"
  visibility = "public"
  page_no = 1
  page_size = 10
}

data "ctyun_ecs_flavors" "ecs_flavor_test" {
  cpu    = 2
  ram    = 4
  arch   = "x86"
  series = "C"
  type   = "CPU_C7"
}

# 创建数据盘资源
resource "ctyun_ebs" "data_disk_test" {
  count      = local.ecs_instance_id == "" ? 1 : 0
  name       = "tf-test-data-disk"
  mode       = "vbd"
  type       = "sata"
  size       = 60
  cycle_type = "on_demand"
}

locals {
  # 查找已存在的数据盘
  data_disks = [for disk in data.ctyun_ebs_volumes.ebs_test.volumes : disk if disk.name == "tf-test-data-disk"]
  data_disk_id = length(local.data_disks) > 0 ? local.data_disks[0].id : ""

  # 确定实际使用的数据盘ID
  real_data_disk_id = local.data_disk_id == "" ? try(ctyun_ebs.data_disk_test[0].id, "") : local.data_disk_id
}

data "ctyun_ebs_volumes" "ebs_test" {
  page_size = 50
}

# 创建ECS实例资源
resource "ctyun_ecs" "ecs_test" {
  count    = local.ecs_instance_id == "" ? 1 : 0
  instance_name      = "tf-test-ecs-for-paas"
  display_name       = "tf-test-init-ecs"
  flavor_id           = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].id
  image_id            = data.ctyun_images.image_test.images[0].id
  system_disk_type   = "sata"
  system_disk_size   = 60
  vpc_id             = local.real_vpc_id
  subnet_id          = local.real_subnet_id
  security_group_ids = [local.real_security_group_id]
  cycle_type         = "on_demand"
}

# 创建EBS与ECS的关联关系（显式挂载）
resource "ctyun_ebs_association_ecs" "data_disk_association" {
  count       = local.ecs_instance_id == "" && local.data_disk_id == "" ? 1 : 0
  instance_id = ctyun_ecs.ecs_test[0].id
  ebs_id      = ctyun_ebs.data_disk_test[0].id
}
