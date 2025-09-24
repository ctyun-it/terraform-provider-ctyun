data "ctyun_vpcs" "vpc_test" {
  page_size = 50
}

locals {
  vpcs = [for vpc in data.ctyun_vpcs.vpc_test.vpcs : vpc if vpc.name == "tf-vpc-for-paas"]
  data_vpc_id = length(local.vpcs) > 0 ? local.vpcs[0].vpc_id : ""
}

resource "ctyun_vpc" "vpc_test" {
  count    = local.data_vpc_id == "" ? 1 : 0
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
  subnets = [for subnet in data.ctyun_subnets.subnet_test.subnets : subnet if subnet.name == "tf-subnet-for-paas"]
  data_subnet_id = length(local.subnets) > 0 ? local.subnets[0].subnet_id : ""
}

resource "ctyun_subnet" "subnet_test" {
  count    = local.data_vpc_id == "" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-subnet-for-paas"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns         = [
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
  security_groups = [for security_group in data.ctyun_security_groups.security_group_test.security_groups : security_group if security_group.name == "tf-sg-for-paas"]
  data_security_group_id = length(local.security_groups) > 0 ? local.security_groups[0].security_group_id : ""
}

resource "ctyun_security_group" "security_group_test" {
  count    = local.data_vpc_id == "" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-sg-for-paas"
  description = "terraform测试使用"
  lifecycle {
    prevent_destroy = true
  }
}

locals {
  real_security_group_id = local.data_security_group_id == "" ? try(ctyun_security_group.security_group_test[0].id, "") : local.data_security_group_id
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
