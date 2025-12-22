data "ctyun_vpcs" "vpc_test" {
  page_size = 50
}

locals {
  vpcs        = [for vpc in data.ctyun_vpcs.vpc_test.vpcs : vpc if vpc.name == "tf-vpc-for-pgsql"]
  data_vpc_id = length(local.vpcs) > 0 ? local.vpcs[0].vpc_id : ""
}

resource "ctyun_vpc" "vpc_test" {
  count       = local.data_vpc_id == "" ? 1 : 0
  name        = "tf-vpc-for-pgsql"
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
    for subnet in data.ctyun_subnets.subnet_test.subnets : subnet if subnet.name == "tf-subnet-for-pgsql"
  ]
  data_subnet_id = length(local.subnets) > 0 ? local.subnets[0].subnet_id : ""
}

resource "ctyun_subnet" "subnet_test" {
  count       = local.data_vpc_id=="" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-subnet-for-pgsql"
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
    for security_group in data.ctyun_security_groups.security_group_test.security_groups :security_group if security_group.name == "tf-sg-for-pgsql"
  ]
  data_security_group_id = length(local.security_groups) > 0 ? local.security_groups[0].security_group_id : ""

  security_groups2 = [
    for security_group in data.ctyun_security_groups.security_group_test.security_groups :security_group if security_group.name == "tf-sg-for-pgsql2"
  ]
  data_security_group_id2 = length(local.security_groups2) > 0 ? local.security_groups2[0].security_group_id : ""

  security_groups3 = [
    for security_group in data.ctyun_security_groups.security_group_test.security_groups :security_group if security_group.name == "tf-sg-for-pgsql3"
  ]
  data_security_group_id3 = length(local.security_groups3) > 0 ? local.security_groups3[0].security_group_id : ""
}

resource "ctyun_security_group" "security_group_test1" {
  count = local.data_vpc_id == "" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-sg-for-pgsql"
  description = "terraform测试使用"
  lifecycle {
    prevent_destroy = false
  }
}
resource "ctyun_security_group" "security_group_test2" {
  count = local.data_vpc_id == "" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-sg-for-pgsql2"
  description = "terraform测试使用2"
  lifecycle {
    prevent_destroy = false
  }
}
resource "ctyun_security_group" "security_group_test3" {
  count = local.data_vpc_id == "" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-sg-for-pgsql3"
  description = "terraform测试使用3"
  lifecycle {
    prevent_destroy = false
  }
}

locals {
  real_security_group_id1 = local.data_security_group_id == "" ? try(ctyun_security_group.security_group_test1[0].id, "") : local.data_security_group_id
  real_security_group_id2 = local.data_security_group_id2 == "" ? try(ctyun_security_group.security_group_test2[0].id, "") : local.data_security_group_id2
  real_security_group_id3 = local.data_security_group_id3 == "" ? try(ctyun_security_group.security_group_test3[0].id, "") : local.data_security_group_id3
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
  flavor_name           = "c7.xlarge.2"
  storage_type          = "SSD"
  storage_space         = 100
  name                  = "pgsql-test-tf2"
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

data "ctyun_postgresql_param_templates" "param_templates" {

}

data "ctyun_postgresql_character_set" "charsets" {

}

data "ctyun_postgresql_collation_time_zone" "collations" {
  depends_on = [ctyun_postgresql_instance.test]
  instance_id    = ctyun_postgresql_instance.test.id
}

resource "ctyun_postgresql_account" "account_test" {
  project_id = "0"
  instance_id = ctyun_postgresql_instance.test.id
  name = "kqjwyk"
  password = var.password
  user_type = "normal"
  description = "terraform测试预置条件"
}

data "ctyun_postgresql_accounts" "accounts" {
  depends_on = [ctyun_postgresql_account.account_test]
  instance_id = ctyun_postgresql_instance.test.id
}

resource "ctyun_postgresql_database" "test" {
  project_id   = "0"
  instance_id      = ctyun_postgresql_instance.test.id
  name         = "test"
  charset_name = "UTF8"
  owner        = ctyun_postgresql_account.account_test.name
}

resource "ctyun_postgresql_database" "test1" {
  project_id   = "0"
  instance_id      = ctyun_postgresql_instance.test.id
  name         = "test1"
  charset_name = "UTF8"
  owner        = ctyun_postgresql_account.account_test.name
  depends_on = [ctyun_postgresql_database.test]
}
