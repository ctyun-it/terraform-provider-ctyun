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
  name        = "vpc-test-mq"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "subnet-test-mq"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "sg-test-mq"
  description = "terraform测试使用"
}

data "ctyun_zones" "test"{

}

data "ctyun_rabbitmq_specs" "test"{

}

locals {
  single_sku = [for sku in data.ctyun_rabbitmq_specs.test.specs[0].sku : sku if sku.prod_name == "单机版"]
  single_disk_type = local.single_sku[0].disk_item.res_items[0]
  single_spec_name = local.single_sku[0].res_item.res_items[0].spec[0].spec_name
}

resource "ctyun_rabbitmq_instance" "test" {
  instance_name = "tf-rabbitmq-example"
  spec_name = local.single_spec_name
  node_num = 1
  zone_list = [data.ctyun_zones.test.zones[0]]
  disk_type = local.single_disk_type
  disk_size = 300
  vpc_id = ctyun_vpc.vpc_test.id
  subnet_id = ctyun_subnet.subnet_test.id
  security_group_id = ctyun_security_group.security_group_test.id
  cycle_type = "on_demand"
}
