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
  name        = "tf-vpc-for-paas"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-paas"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-paas"
  description = "terraform测试使用"
  lifecycle {
    prevent_destroy = true
  }
}

data "ctyun_mysql_specs" "mysql_specs" {
  instance_series = "S"
}

// mysql创建单节点
resource "ctyun_mysql_instance" "mysql_test" {
  cycle_type            = "on_demand"
  vpc_id                = ctyun_vpc.vpc_test.id
  flavor_name           = "c7.2xlarge.4"
  subnet_id             = ctyun_subnet.subnet_test.id
  security_group_id     = ctyun_security_group.security_group_test.id
  name                  = "mysql_examples"
  prod_id               = "Single57"
  storage_type          = "SATA"
  storage_space         = 100
}

// mysql创建1主1备
resource "ctyun_mysql_instance" "mysql_test" {
  cycle_type            = "on_demand"
  vpc_id                = ctyun_vpc.vpc_test.id
  flavor_name           = "c7.2xlarge.4"
  subnet_id             = ctyun_subnet.subnet_test.id
  security_group_id     = ctyun_security_group.security_group_test.id
  name                  = "mysql_examples"
  prod_id               = "MasterSlave80"
  storage_type          = "SATA"
  storage_space         = 100
}
// 升配磁盘空间

resource "ctyun_mysql_instance" "mysql_test" {
  cycle_type            = "on_demand"
  vpc_id                = ctyun_vpc.vpc_test.id
  flavor_name           = "c7.2xlarge.4"
  subnet_id             = ctyun_subnet.subnet_test.id
  security_group_id     = ctyun_security_group.security_group_test.id
  name                  = "mysql_examples"
  prod_id               = "Single57"
  storage_type          = "SATA"
  storage_space         = 120
  availability_zone_info = [
    { "availability_zone_name" : "cn-gs-qyi2-1a-public-ctcloud", "availability_zone_count" : 1, "node_type" : "master" }
  ]
}

// 升配规格

resource "ctyun_mysql_instance" "mysql_test" {
  cycle_type            = "on_demand"
  vpc_id                = ctyun_vpc.vpc_test.id
  flavor_name           = "c7.2xlarge.4"
  subnet_id             = ctyun_subnet.subnet_test.id
  security_group_id     = ctyun_security_group.security_group_test.id
  name                  = "mysql_examples"
  prod_id               = "Single57"
  storage_type          = "SATA"
  storage_space         = 120
  availability_zone_info = [
    { "availability_zone_name" : "cn-gs-qyi2-1a-public-ctcloud", "availability_zone_count" : 1, "node_type" : "master" }
  ]
}

// 升配节点 (如单节点->一主两备)
resource "ctyun_mysql_instance" "mysql_test" {
  cycle_type            = "on_demand"
  vpc_id                = ctyun_vpc.vpc_test.id
  flavor_name           = "c7.2xlarge.4"
  subnet_id             = ctyun_subnet.subnet_test.id
  security_group_id     = ctyun_security_group.security_group_test.id
  name                  = "mysql_examples"
  prod_id               = "Master2Slave57"
  storage_type          = "SATA"
  storage_space         = 120
  availability_zone_info = [
    { "availability_zone_name" : "cn-gs-qyi2-1a-public-ctcloud", "availability_zone_count" : 2, "node_type" : "slave" } // 当升配时，availability_zone_info仅需要填写增量的节点信息
  ]
}