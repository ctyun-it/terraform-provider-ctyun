terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 可参考index.md，在环境变量中配置ak、sk、资源池ID、可用区名称
provider "ctyun" {
  env = "prod"
}


resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-mysql"
  cidr        = "192.168.0.0/16"
  description = "terraform-mysql测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-mysql1"
  cidr        = "192.168.1.0/24"
  description = "terraform-mysql测试使用"
  dns = [
    "114.114.114.114",
    "8.8.8.8",
  ]
}
resource "ctyun_security_group" "sg_mysql_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-mysql"
  description = "terraform-mysql测试使用"
  lifecycle {
    prevent_destroy = false
  }
}

resource "ctyun_mysql_instance" "mysql_src" {
  cycle_type        = "on_demand"
  vpc_id            = ctyun_vpc.vpc_test.id
  subnet_id         = ctyun_subnet.subnet_test.id
  security_group_id = ctyun_security_group.sg_mysql_test.id
  name              = "mysql-test-backup-1"
  prod_id           = "Single57"
  storage_type      = "SATA"
  storage_space     = 100
  password          = var.password
  flavor_name       = "c7.2xlarge.4"
}


variable "password" {
  type      = string
  sensitive = true
}
resource "ctyun_mysql_database" "db1" {
  instance_id  = ctyun_mysql_instance.mysql_src.id
  name         = "test_db1"
  charset_name = "utf8mb4"
}
resource "ctyun_mysql_database" "db2" {
  instance_id  = ctyun_mysql_instance.mysql_src.id
  name         = "test_db2"
  charset_name = "utf8mb4"
  depends_on   = [ctyun_mysql_database.db1]
}
resource "ctyun_mysql_database" "db3" {
  instance_id  = ctyun_mysql_instance.mysql_src.id
  name         = "test_db3"
  charset_name = "utf8mb4"
  depends_on   = [ctyun_mysql_database.db2]
}

resource "ctyun_mysql_backup" "backup_test" {
  instance_id = ctyun_mysql_instance.mysql_src.id
  project_id  = "0"
  description = "terraform单元测试"
  task_type   = "full"
  depends_on  = [ctyun_mysql_database.db3]
}

data "ctyun_mysql_backups" "example" {
  instance_id = ctyun_mysql_instance.mysql_src.id
}


resource "ctyun_mysql_backup_cancel" "example" {
  instance_id      = ctyun_mysql_instance.mysql_src.id
  backup_record_id = data.ctyun_mysql_backups.example.backups.0.records.0.backup_record_id
}
