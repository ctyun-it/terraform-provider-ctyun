data "ctyun_vpcs" "vpc_test" {
  page_size = 50
}

locals {
  vpcs        = [for vpc in data.ctyun_vpcs.vpc_test.vpcs : vpc if vpc.name == "tf-vpc-for-paas"]
  data_vpc_id = length(local.vpcs) > 0 ? local.vpcs[0].vpc_id : ""
}

resource "ctyun_vpc" "vpc_test" {
  count       = local.data_vpc_id == "" ? 1 : 0
  name        = "tf-vpc-for-paas"
  cidr        = "192.168.0.0/16"
  description = "terraform-paas测试使用"
  enable_ipv6 = true
}

locals {
  real_vpc_id = local.data_vpc_id == "" ? try(ctyun_vpc.vpc_test[0].id, "") : local.data_vpc_id
}


data "ctyun_subnets" "subnet_test" {
  vpc_id = local.real_vpc_id
}

locals {
  subnets = [
    for subnet in data.ctyun_subnets.subnet_test.subnets : subnet if subnet.name == "tf-subnet-for-paas"
  ]
  data_subnet_id = length(local.subnets) > 0 ? local.subnets[0].subnet_id : ""
}

resource "ctyun_subnet" "subnet_test" {
  count       = local.data_vpc_id=="" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-subnet-for-paas"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns = [
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
  security_groups = [
    for security_group in data.ctyun_security_groups.security_group_test.security_groups :security_group if security_group.name == "tf-sg-for-paas"
  ]
  data_security_group_id = length(local.security_groups) > 0 ? local.security_groups[0].security_group_id : ""

  security_groups2 = [
    for security_group in data.ctyun_security_groups.security_group_test.security_groups :security_group if security_group.name == "tf-sg-for-paas2"
  ]
  data_security_group_id2 = length(local.security_groups2) > 0 ? local.security_groups2[0].security_group_id : ""
}

resource "ctyun_security_group" "security_group_test1" {
  count = local.data_vpc_id=="" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-sg-for-paas"
  description = "terraform测试使用"
  lifecycle {
    prevent_destroy = false
  }
}
resource "ctyun_security_group" "security_group_test2" {
  count = local.data_vpc_id=="" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-sg-for-paas2"
  description = "terraform测试使用2"
  lifecycle {
    prevent_destroy = false
  }
}

locals {
  real_security_group_id1 = local.data_security_group_id == "" ? try(ctyun_security_group.security_group_test1[0].id, "") : local.data_security_group_id
  real_security_group_id2 = local.data_security_group_id2 == "" ? try(ctyun_security_group.security_group_test2[0].id, "") : local.data_security_group_id
}

resource "ctyun_eip" "eip_test" {
  name                = "tf-eip-for-pgsql"
  bandwidth           = 1
  cycle_type          = "on_demand"
  demand_billing_type = "upflowc"
}



resource "ctyun_postgresql_instance" "test" {
  cycle_type            = "on_demand"
  prod_id               = "Single1222"
  flavor_name           = "s7.large.2"
  storage_type          = "SATA"
  storage_space         = 100
  name                  = "pgsql-test-2"
  password              = var.password
  case_sensitive        = true
  vpc_id                = local.real_vpc_id
  subnet_id             = local.real_subnet_id
  security_group_id     = local.real_security_group_id1
  backup_storage_type  = "OS"
}

variable "password" {
  type      = string
  sensitive = true
}

data "ctyun_zones" "az" {

}
