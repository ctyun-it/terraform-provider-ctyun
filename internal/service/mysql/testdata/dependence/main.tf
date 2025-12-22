// main.tf负责创建或查询单测依赖的前置资源
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
  subnets = [
    for subnet in data.ctyun_subnets.subnet_test.subnets : subnet if subnet.name == "tf-subnet-for-paas"
  ]
  data_subnet_id = length(local.subnets) > 0 ? local.subnets[0].subnet_id : ""
}

resource "ctyun_subnet" "subnet_test" {
  count       = local.data_vpc_id=="" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-subnet-for-paas"
  cidr        = "192.168.0.0/16"
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
}

resource "ctyun_security_group" "security_group_test" {
  count       = local.data_vpc_id=="" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-sg-for-paas"
  description = "terraform测试使用"
  lifecycle {
    prevent_destroy = false
  }
}

locals {
  real_security_group_id = local.data_security_group_id == "" ? try(ctyun_security_group.security_group_test[0].id, "") : local.data_security_group_id
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

data "ctyun_mysql_specs" "mysql_specs"{
  instance_series = "S"
}

data "ctyun_mysql_backups" "backup_test" {
  depends_on = [ctyun_mysql_backup.backup_test]
  instance_id   = ctyun_mysql_instance.mysql_test.id
  page_no   = 1
  page_size = 10
}

resource "ctyun_mysql_instance" "mysql_test" {
  cycle_type            = "on_demand"
  vpc_id                = local.real_vpc_id
  flavor_name           = "c7.large.2"
  prod_id               = "Single57"
  subnet_id             = local.real_subnet_id
  security_group_id     = local.real_security_group_id
  name                  = local.mysql_name
  storage_type          = "SATA"
  storage_space         = 100
  lifecycle {
    ignore_changes = [name]
  }
}

resource "ctyun_mysql_backup" "backup_test" {
  instance_id     = ctyun_mysql_instance.mysql_test.id
  project_id  = "0"
  description = "terraform单元测试"
  task_type   = "full"
  depends_on = [ctyun_mysql_database.db3]
}

data "ctyun_mysql_recoverable_time_points" "time_point_test" {
  depends_on = [ctyun_mysql_backup.backup_test]
  instance_id    = ctyun_mysql_instance.mysql_test.id
  project_id = "0"
}

data "ctyun_mysql_param_templates" "template"{
  engine = "5.7"
  name = "parameterSet57"
}

locals {
  # 生成当前时间戳的哈希值
  hash = sha256(timestamp())

  # 从哈希结果中截取字符（转为小写并移除特殊字符）
  random_string = substr(
    replace(
      lower(local.hash),
      "/[^a-z0-9]/",
      ""  # 移除所有非字母数字的字符
    ),
    0, 5  # 截取前10个字符
  )
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