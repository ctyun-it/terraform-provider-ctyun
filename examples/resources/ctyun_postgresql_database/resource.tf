terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

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

data "ctyun_postgresql_character_set" "charsets" {

}

data "ctyun_postgresql_collation_time_zone" "collations" {
  depends_on  = [ctyun_postgresql_instance.test]
  instance_id = ctyun_postgresql_instance.test.id
}

resource "ctyun_postgresql_database" "examples" {
  project_id      = "0"
  instance_id     = ctyun_postgresql_instance.test.id
  name            = "pg_test"
  charset_name    = data.ctyun_postgresql_character_set.charsets.character_set[1]
  charset_collate = data.ctyun_postgresql_collation_time_zone.collations.collations[0].coll_name
  charset_type    = data.ctyun_postgresql_collation_time_zone.collations.collations[0].coll_type
  owner           = "root"
  description     = "postgresql 样例"
}
