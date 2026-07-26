resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-redis"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-redis"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns         = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-redis"
  description = "terraform测试使用"
}

resource "ctyun_security_group_rule" "security_group_rule_ingress" {
 security_group_id = ctyun_security_group.security_group_test.id
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
  vpc_id = ctyun_vpc.vpc_test.id
  subnet_id = ctyun_subnet.subnet_test.id
  security_group_id = ctyun_security_group.security_group_test.id
  password=var.password
  cycle_type = "month"
  cycle_count = 1
  shard_mem_size = 8
  host_type = "C"
}

resource "ctyun_redis_instance" "test_redis_instance2" {
  instance_name = "test-redis-instance6"
  engine_version = "7.0"
  edition = local.spec.series_code
  vpc_id = ctyun_vpc.vpc_test.id
  subnet_id = ctyun_subnet.subnet_test.id
  security_group_id = ctyun_security_group.security_group_test.id
  password=var.password
  cycle_type = "month"
  cycle_count = 1
  auto_renew = false
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