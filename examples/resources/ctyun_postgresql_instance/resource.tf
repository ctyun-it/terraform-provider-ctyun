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
    "8.8.4.4"
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
  cycle_type            = "on_demand"
  host_type             = "S7"
  prod_id               = "Single1222"
  storage_type          = "SATA"
  storage_space         = 100
  name                  = "pgsql-test"
  password              = var.password
  case_sensitive        = true
  instance_series       = "S"
  prod_performance_spec = "2C4G"
  vpc_id                = ctyun_vpc.vpc_test.id
  subnet_id             = ctyun_subnet.subnet_test.id
  security_group_id     = ctyun_security_group.sg_test.id
  availability_zone_info = [
    { "availability_zone_name" : "cn-gs-qyi2-1a-public-ctcloud", "availability_zone_count" : 1, "node_type" : "master" }
  ] // availability_zone_name值根据情况而定
  backup_storage_type   = "SATA"
  backup_storage_space  = 100
  os_type               = "ctyunos"
  cpu_type              = "Intel"
}

// 升配pgsql--对备用磁盘扩容(在升配主storage时候，确保备用磁盘空间>主磁盘空间)
resource "ctyun_postgresql_instance" "test" {
  cycle_type            = "on_demand"
  host_type             = "S7"
  prod_id               = "Single1222"
  storage_type          = "SATA"
  storage_space         = 100
  name                  = "pgsql-test"
  password              = var.password
  case_sensitive        = true
  instance_series       = "S"
  prod_performance_spec = "2C4G"
  vpc_id                = ctyun_vpc.vpc_test.id
  subnet_id             = ctyun_subnet.subnet_test.id
  security_group_id     = ctyun_security_group.sg_test.id
  availability_zone_info = [
    { "availability_zone_name" : "cn-gs-qyi2-1a-public-ctcloud", "availability_zone_count" : 1, "node_type" : "master" }
  ] // availability_zone_name值根据情况而定
  backup_storage_type   = "SATA"
  backup_storage_space  = 120
  os_type               = "ctyunos"
  cpu_type              = "Intel"
}


// 升配pgsql--对主磁盘扩容(在升配主storage时候，确保备用磁盘空间>主磁盘空间)
resource "ctyun_postgresql_instance" "test" {
  cycle_type            = "on_demand"
  host_type             = "S7"
  prod_id               = "Single1222"
  storage_type          = "SATA"
  storage_space         = 120
  name                  = "pgsql-test"
  password              = var.password
  case_sensitive        = true
  instance_series       = "S"
  prod_performance_spec = "2C4G"
  vpc_id                = ctyun_vpc.vpc_test.id
  subnet_id             = ctyun_subnet.subnet_test.id
  security_group_id     = ctyun_security_group.sg_test.id
  availability_zone_info = [
    { "availability_zone_name" : "cn-gs-qyi2-1a-public-ctcloud", "availability_zone_count" : 1, "node_type" : "master" }
  ] // availability_zone_name值根据情况而定
  backup_storage_type   = "SATA"
  backup_storage_space  = 120
  os_type               = "ctyunos"
  cpu_type              = "Intel"
}


// 升配规格 2C4G->2C8G
resource "ctyun_postgresql_instance" "test" {
  cycle_type            = "on_demand"
  host_type             = "S7"
  prod_id               = "Single1222"
  storage_type          = "SATA"
  storage_space         = 120
  name                  = "pgsql-test"
  password              = var.password
  case_sensitive        = true
  instance_series       = "S"
  prod_performance_spec = "2C8G"
  vpc_id                = ctyun_vpc.vpc_test.id
  subnet_id             = ctyun_subnet.subnet_test.id
  security_group_id     = ctyun_security_group.sg_test.id
  availability_zone_info = [
    { "availability_zone_name" : "cn-gs-qyi2-1a-public-ctcloud", "availability_zone_count" : 1, "node_type" : "master" }
  ] // availability_zone_name值根据情况而定
  backup_storage_type   = "SATA"
  backup_storage_space  = 120
  os_type               = "ctyunos"
  cpu_type              = "Intel"
}


// 升配类型

// 升配规格 单节点->1主2备
resource "ctyun_postgresql_instance" "test" {
  cycle_type            = "on_demand"
  host_type             = "S7"
  prod_id               = "Master2Slave1222"
  storage_type          = "SATA"
  storage_space         = 120
  name                  = "pgsql-test"
  password              = var.password
  case_sensitive        = true
  instance_series       = "S"
  prod_performance_spec = "2C8G"
  vpc_id                = ctyun_vpc.vpc_test.id
  subnet_id             = ctyun_subnet.subnet_test.id
  security_group_id     = ctyun_security_group.sg_test.id
  availability_zone_info = [
    { "availability_zone_name" : "cn-gs-qyi2-1a-public-ctcloud", "availability_zone_count" : 2, "node_type" : "slave" }
  ] // availability_zone_name值根据情况而定
  backup_storage_type   = "SATA"
  backup_storage_space  = 120
  os_type               = "ctyunos"
  cpu_type              = "Intel"
}



