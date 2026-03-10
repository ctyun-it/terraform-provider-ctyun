resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-rabbitmq"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-rabbitmq"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns         = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-rabbitmq"
  description = "terraform测试使用"
}

data "ctyun_rabbitmq_specs" "test"{

}

locals {
  single_sku = [for sku in data.ctyun_rabbitmq_specs.test.specs[0].sku : sku if sku.prod_name == "单机版"]
  single_disk_type = local.single_sku[0].disk_item.res_items[0]
  single_spec_name = local.single_sku[0].res_item.res_items[0].spec[0].spec_name
  single_spec_name2 = local.single_sku[0].res_item.res_items[0].spec[1].spec_name

  cluster_sku = [for sku in data.ctyun_rabbitmq_specs.test.specs[0].sku : sku if sku.prod_name == "集群版"]
  cluster_disk_type = local.cluster_sku[0].disk_item.res_items[0]
  cluster_spec_name = local.cluster_sku[0].res_item.res_items[0].spec[0].spec_name
  cluster_spec_name2 = local.cluster_sku[0].res_item.res_items[0].spec[1].spec_name
}

data "ctyun_zones" "test" {

}

resource "ctyun_rabbitmq_instance" "test" {
  instance_name = "tf-rabbitmq-${local.random_string}"
  spec_name = local.cluster_spec_name
  node_num = 3
  zone_list = [data.ctyun_zones.test.zones[0]]
  disk_type = local.cluster_disk_type
  disk_size = 300
  vpc_id = ctyun_vpc.vpc_test.id
  subnet_id = ctyun_subnet.subnet_test.id
  security_group_id = ctyun_security_group.security_group_test.id
  cycle_type = "on_demand"
}

resource "ctyun_rabbitmq_exchange" "test" {
  instance_id = ctyun_rabbitmq_instance.test.id
  vhost = "/"
  name = "tf-exchange"
  type = "direct"
}

locals {
  # 生成当前时间戳的哈希值
  hash = sha256(timestamp())

  # 从哈希结果中截取字符（转为小写并移除特殊字符）
  random_string = substr(
    replace(
      lower(local.hash),
      "/[^a-z0-9]/",
      ""  # 移除所有非字母数字的字符
    ),
    0, 10  # 截取前16个字符
  )
}
