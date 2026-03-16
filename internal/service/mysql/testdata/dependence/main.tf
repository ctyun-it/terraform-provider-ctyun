resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-mysql"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-mysql"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-mysql"
  description = "terraform测试使用"
}

resource "ctyun_eip" "eip_test" {
  name                = "tf-eip-for-mysql"
  bandwidth           = 1
  cycle_type          = "on_demand"
  demand_billing_type = "upflowc"
}

data "ctyun_zones" "test" {

}

locals {
  mysql_name = "tf-mysql-${local.random_string}"
  az_name    = data.ctyun_zones.test.zones[0]
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

data "ctyun_mysql_backups" "backup_test" {
  depends_on = [ctyun_mysql_backup.backup_test]
  instance_id   = ctyun_mysql_instance.mysql_test.id
  page_no   = 1
  page_size = 10
}

resource "ctyun_mysql_instance" "mysql_test" {
  cycle_type            = "on_demand"
  vpc_id                = ctyun_vpc.vpc_test.id
  flavor_name           = local.flavor_name
  prod_id               = "Single57"
  subnet_id             = ctyun_subnet.subnet_test.id
  security_group_id     = ctyun_security_group.security_group_test.id
  name                  = local.mysql_name
  storage_type          = "SSD"
  storage_space         = 100
  lifecycle {
    ignore_changes = [name]
  }
}

resource "ctyun_mysql_backup" "backup_test" {
  instance_id     = ctyun_mysql_instance.mysql_test.id
  description = "terraform单元测试"
  task_type   = "full"
  depends_on = [ctyun_mysql_database.db3]
}

data "ctyun_mysql_recoverable_time_points" "time_point_test" {
  depends_on = [ctyun_mysql_backup.backup_test]
  instance_id    = ctyun_mysql_instance.mysql_test.id
}

data "ctyun_mysql_param_templates" "template"{
  engine = "5.7"
  name = "parameterSet57"
}

locals {
  # 生成当前时间戳的哈希值
  random_string = substr(replace(lower(sha256(timestamp())), "/[^a-z0-9]/", ""), 0, 5)
}

resource "ctyun_mysql_database" "db1" {
  instance_id      = ctyun_mysql_instance.mysql_test.id
  name         = "test_db1"
  charset_name = "utf8mb4"
}
resource "ctyun_mysql_database" "db2" {
  instance_id      = ctyun_mysql_instance.mysql_test.id
  name         = "test_db2"
  charset_name = "utf8mb4"
  depends_on = [ctyun_mysql_database.db1]
}
resource "ctyun_mysql_database" "db3" {
  instance_id  = ctyun_mysql_instance.mysql_test.id
  name         = "test_db3"
  charset_name = "utf8mb4"
  depends_on = [ctyun_mysql_database.db2]
}