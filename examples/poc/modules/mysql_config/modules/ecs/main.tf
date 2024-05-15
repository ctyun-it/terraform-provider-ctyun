terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
}

data "ctyun_ecs_flavor" "ecs_flavor_test" {
  cpu    = 1
  ram    = 1
  arch   = "x86"
  series = "S"
}

resource "ctyun_ecs" "mysql_test" {
  count              = var.ecs_instance_count
  name               = "${var.ecs_name}-${count.index}"
  flavor_id          = data.ctyun_ecs_flavor.ecs_flavor_test.flavors[0].id
  image_id           = "9a099800-3e1c-45cd-99d1-7e2207a2fb08"
  system_disk_type   = "SATA"
  system_disk_size   = 40
  vpc_id             = var.vpc_id
  security_group_ids = var.security_group_ids
  password           = var.ecs_password
  cycle_type         = "month"
  cycle_count        = 1
  subnet_id          = var.subnet_id
  auto_renew         = true
  status             = "running"
  key_pair_name      = var.ecs_key_pair_name
}