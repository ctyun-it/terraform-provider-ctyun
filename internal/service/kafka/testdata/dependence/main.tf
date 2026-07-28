resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-kafka"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-kafka"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns         = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-kafka"
  description = "terraform测试使用"
}

data "ctyun_kafka_specs" "test"{

}

locals {
  single_sku = [for sku in data.ctyun_kafka_specs.test.specs[0].sku : sku if sku.prod_name == "单机版"]
  single_disk_type = local.single_sku[0].disk_item.res_items[0]
  single_spec_name = local.single_sku[0].res_item.res_items[0].spec[0].spec_name
  single_spec_name2 = local.single_sku[0].res_item.res_items[0].spec[1].spec_name

  cluster_sku = [for sku in data.ctyun_kafka_specs.test.specs[0].sku : sku if sku.prod_name == "集群版"]
  cluster_disk_type = local.cluster_sku[0].disk_item.res_items[0]
  cluster_spec_name = local.cluster_sku[0].res_item.res_items[0].spec[0].spec_name
  cluster_spec_name2 = local.cluster_sku[0].res_item.res_items[0].spec[1].spec_name

}

data "ctyun_zones" "test" {

}

resource "ctyun_kafka_instance" "test_kafka_instance" {
  instance_name = "tf-kafka-inst1"
  engine_version = "3.6"
  spec_name = local.cluster_spec_name
  node_num = 3
  zone_list = [data.ctyun_zones.test.zones[0]]
  disk_type = local.cluster_disk_type
  disk_size = 100
  vpc_id = ctyun_vpc.vpc_test.id
  subnet_id = ctyun_subnet.subnet_test.id
  security_group_id = ctyun_security_group.security_group_test.id
  retention_hours = 80
  cycle_type = "month"
  cycle_count = 2
  auto_renew = true
  auto_renew_cycle_count = 1
}

resource "ctyun_kafka_topic" "test_kafka_topic" {
  name = "test_topic"
  instance_id = ctyun_kafka_instance.test_kafka_instance.id
  partition_num  = 1
}

resource "ctyun_kafka_user" "test_kafka_user" {
  name = "test_kafka_user"
  instance_id = ctyun_kafka_instance.test_kafka_instance.id
  password  = var.password
}

variable "password" {
  type      = string
  sensitive = true
}
