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
  name        = "tf-sg-for-esc"
  description = "terraform-mysql测试使用"
  lifecycle {
    prevent_destroy = false
  }
}

resource "ctyun_mysql_instance" "mysql_test" {
  cycle_type        = "on_demand"
  vpc_id            = ctyun_vpc.vpc_test.id
  subnet_id         = ctyun_subnet.subnet_test.id
  security_group_id = ctyun_security_group.sg_mysql_test.id
  name              = "mysql-test-1"
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

# 配置文件: resource_ctyun_mysql_audit.tf
resource "ctyun_mysql_audit" "example" {
  instance_id  = ctyun_mysql_instance.mysql_test.id
  project_id   = "0"
  audit_switch = true
}