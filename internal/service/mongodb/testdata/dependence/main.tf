resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-mongodb"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-mongodb"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-mongodb"
  description = "terraform测试使用"
}

resource "ctyun_eip" "eip_test" {
  name                = "tf-eip-for-mongodb"
  bandwidth           = 1
  cycle_type          = "on_demand"
  demand_billing_type = "upflowc"
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

resource "ctyun_mongodb_instance" "mongodb_eip" {
  cycle_type             = "on_demand"
  vpc_id                 = ctyun_vpc.vpc_test.id
  flavor_name            = local.flavor_name
  subnet_id              = ctyun_subnet.subnet_test.id
  security_group_id      =  ctyun_security_group.security_group_test.id
  name                   = "mongodb-${local.random_string}"
  prod_id                = "Single34"
  storage_type           = "SATA"
  storage_space          = 100
  backup_storage_type    = "OS"
  password = var.password
  lifecycle {
    ignore_changes = [name]
  }
}

resource "ctyun_mongodb_instance" "mongodb_ei2p" {
  cycle_type             = "month"
  cycle_count  = 2
  vpc_id                 = ctyun_vpc.vpc_test.id
  flavor_name            = local.flavor_name
  subnet_id              = ctyun_subnet.subnet_test.id
  security_group_id      =  ctyun_security_group.security_group_test.id
  name                   = "mongodb-${local.random_string}22"
  prod_id                = "Single34"
  storage_type           = "SATA"
  storage_space          = 100
  backup_storage_type    = "OS"
  password = var.password
  lifecycle {
    ignore_changes = [name]
  }
}

variable "password" {
  type      = string
  sensitive = true
}

locals {
  random_string = substr(replace(lower(sha256(timestamp())), "/[^a-z0-9]/", ""), 0, 5)
}

data "ctyun_zones" "az" {

}