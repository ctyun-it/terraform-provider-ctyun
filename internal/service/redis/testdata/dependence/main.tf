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

resource "ctyun_security_group_rule" "security_group_rule_ingress" {
 security_group_id = local.real_security_group_id
 direction         = "ingress"
 action            = "accept"
 priority          = 1
 protocol          = "tcp"
 ether_type        = "ipv4"
 dest_cidr_ip      = "0.0.0.0/0"
 range             = "6379"
}


resource "ctyun_eip" "eip_test" {
  name                = "tf-eip-for-redis"
  bandwidth           = 10
  cycle_type          = "on_demand"
  demand_billing_type = "bandwidth"
}

data "ctyun_redis_specs" "test"{

}

locals {
  spec = data.ctyun_redis_specs.test.series_infos[0]
}

resource "ctyun_redis_instance" "test_redis_instance" {
  instance_name = "test-redis-instance7"
  engine_version = "7.0"
  edition = local.spec.series_code
  vpc_id = local.real_vpc_id
  subnet_id = local.real_subnet_id
  security_group_id = local.real_security_group_id
  password=var.password
  cycle_type = "month"
  cycle_count = 1
  auto_renew = true
  auto_renew_cycle_count = 12
  shard_mem_size = 8
  host_type = "C"
}

resource "ctyun_redis_instance" "test_redis_instance2" {
  instance_name = "test-redis-instance6"
  engine_version = "7.0"
  edition = local.spec.series_code
  vpc_id = local.real_vpc_id
  subnet_id = local.real_subnet_id
  security_group_id = local.real_security_group_id
  password=var.password
  cycle_type = "month"
  cycle_count = 1
  auto_renew = true
  auto_renew_cycle_count = 12
  shard_mem_size = 8
  host_type = "C"
}

resource "ctyun_redis_account" "test_instance1_account" {
  name = "instance1_account"
  instance_id = ctyun_redis_instance.test_redis_instance.id
  password  = var.password
  privilege = "rw"
}

resource "ctyun_redis_account" "test_instance2_account" {
  name = "instance2_account"
  instance_id = ctyun_redis_instance.test_redis_instance2.id
  password  = var.password
  privilege = "rw"
}

variable "password" {
  type      = string
  sensitive = true
}