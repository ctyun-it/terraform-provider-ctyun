

resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-pgsql"
  cidr        = "192.168.0.0/16"
  description = "terraform-paas测试使用"
  enable_ipv6 = true
}


resource "ctyun_subnet" "subnet_test" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-pgsql"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}


resource "ctyun_security_group" "security_group_test1" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-pgsql"
  description = "terraform测试使用"
  lifecycle {
    prevent_destroy = false
  }
}
resource "ctyun_security_group" "security_group_test2" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-pgsql2"
  description = "terraform测试使用2"
  lifecycle {
    prevent_destroy = false
  }
}
resource "ctyun_security_group" "security_group_test3" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-pgsql3"
  description = "terraform测试使用3"
  lifecycle {
    prevent_destroy = false
  }
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
  storage_type          = "SSD"
  storage_space         = 100
  name                  = "pgsql-test-tf4"
  password              = var.password
  case_sensitive        = true
  vpc_id                = ctyun_vpc.vpc_test.id
  subnet_id             = ctyun_subnet.subnet_test.id
  security_group_id     = ctyun_security_group.security_group_test1.id
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
  instance_id      = ctyun_postgresql_instance.test.id
  name         = "test"
  charset_name = "UTF8"
  owner        = ctyun_postgresql_account.account_test.name
}

resource "ctyun_postgresql_database" "test1" {
  instance_id      = ctyun_postgresql_instance.test.id
  name         = "test1"
  charset_name = "UTF8"
  owner        = ctyun_postgresql_account.account_test.name
  depends_on = [ctyun_postgresql_database.test]
}



data "ctyun_ecs_flavors" "ecs_flavor_test" {
  cpu    = 2
  ram    = 4
  arch   = "x86"
}

data "ctyun_ecs_flavors" "ecs_flavor_test2" {
  cpu    = 4
  ram    = 8
  arch   = "x86"
}

locals {
  flavor_name = [for f in data.ctyun_ecs_flavors.ecs_flavor_test.flavors : f.name if f.available == true][0]
  flavor_name2 = [for f in data.ctyun_ecs_flavors.ecs_flavor_test2.flavors : f.name if f.available == true][0]
}