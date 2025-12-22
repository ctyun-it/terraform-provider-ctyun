

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
resource "ctyun_security_group" "sg_pgsql_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-esc"
  description = "terraform-kafka测试使用"
  lifecycle {
    prevent_destroy = false
  }
}

variable "password" {
  type      = string
  sensitive = true
}

resource "ctyun_postgresql_instance" "test" {
  cycle_type          = "on_demand"
  prod_id             = "Single1222"
  flavor_name         = "c7.xlarge.4"
  storage_type        = "SSD"
  storage_space       = 120
  name                = "pgsql-test-tf1"
  password            = var.password
  case_sensitive      = true
  vpc_id              = ctyun_vpc.vpc_test.id
  subnet_id           = ctyun_subnet.subnet_test.id
  security_group_id   = ctyun_security_group.sg_pgsql_test.id
  backup_storage_type = "OS"
}

resource "ctyun_eip" "eip_test" {
  name                = "tf-eip-for-pgsql"
  bandwidth           = 1
  cycle_type          = "on_demand"
  demand_billing_type = "upflowc"
}

resource "ctyun_postgresql_association_eip" "pgsql_association_eip_test" {
  eip_id      = ctyun_eip.eip_test.id
  instance_id = ctyun_postgresql_instance.test.id
}
