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

variable "password" {
  type      = string
  sensitive = true
}

resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-pgsql"
  cidr        = "192.168.0.0/16"
  description = "terraform-kafka测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-pgsql"
  cidr        = "192.168.1.0/24"
  description = "terraform-kafka测试使用"
  dns = [
    "114.114.114.114",
    "8.8.8.8",
  ]
}
resource "ctyun_security_group" "sg_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-esc"
  description = "terraform-kafka测试使用"
  lifecycle {
    prevent_destroy = false
  }
}

// 开通样例
resource "ctyun_postgresql_instance" "test" {
  cycle_type          = "on_demand"
  prod_id             = "Single1222"
  flavor_name         = "c7.xlarge.2"
  storage_type        = "SSD"
  storage_space       = 100
  name                = "pgsql-test-tf1"
  password            = var.password
  case_sensitive      = true
  vpc_id              = ctyun_vpc.vpc_test.id
  subnet_id           = ctyun_subnet.subnet_test.id
  security_group_id   = ctyun_security_group.sg_test.id
  backup_storage_type = "OS"
}


resource "ctyun_postgresql_account" "account_test" {
  instance_id = ctyun_postgresql_instance.test.id
  name        = "example_user"
  password    = var.password
  user_type   = "normal"
  description = "terraform样例"
}