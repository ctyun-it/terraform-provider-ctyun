data "ctyun_vpcs" "vpc_test" {
  page_size = 50
}

locals {
  vpcs = [for vpc in data.ctyun_vpcs.vpc_test.vpcs : vpc if vpc.name == "tf-vpc-for-rocketmq"]
  data_vpc_id = length(local.vpcs) > 0 ? local.vpcs[0].vpc_id : ""
}

resource "ctyun_vpc" "vpc_test" {
  count    = local.data_vpc_id == "" ? 1 : 0
  name        = "tf-vpc-for-rocketmq"
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
  subnets = [for subnet in data.ctyun_subnets.subnet_test.subnets : subnet if subnet.name == "tf-subnet-for-rocketmq"]
  data_subnet_id = length(local.subnets) > 0 ? local.subnets[0].subnet_id : ""
}

resource "ctyun_subnet" "subnet_test" {
  count    = local.data_vpc_id == "" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-subnet-for-rocketmq"
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
  security_groups = [for security_group in data.ctyun_security_groups.security_group_test.security_groups : security_group if security_group.name == "tf-sg-for-rocketmq"]
  data_security_group_id = length(local.security_groups) > 0 ? local.security_groups[0].security_group_id : ""
}

resource "ctyun_security_group" "security_group_test" {
  count    = local.data_vpc_id == "" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-sg-for-rocketmq"
  description = "terraform测试使用"
}

locals {
  real_security_group_id = local.data_security_group_id == "" ? try(ctyun_security_group.security_group_test[0].id, "") : local.data_security_group_id
}

data "ctyun_rocketmq_specs" "test"{

}

locals {
  # 筛选出 spec_name 以 "single" 结尾的规格（单机版）
  single_specs = [for spec in data.ctyun_rocketmq_specs.test.specs : spec if endswith(spec.spec_name, "single")]
  # 取前两个单机版规格
  single_spec_name  ="rocketmq.2u4g.single"
  single_spec_name2 = "rocketmq.4u8g.single"

  # 筛选出 spec_name 不以 "single" 结尾的规格（集群版等）
  cluster_specs = [for spec in data.ctyun_rocketmq_specs.test.specs : spec if !endswith(spec.spec_name, "single")]
  # 取前两个非单机版规格
  cluster_spec_name  = "rocketmq.4u8g.cluster"
  cluster_spec_name2 = "rocketmq.8u16g.cluster"
}

data "ctyun_zones" "test" {
}

resource"ctyun_rocketmq_instance" "example" {
  instance_name = "tf-rocketmq-instance-${local.random_string}"
  spec_name = "rocketmq.4u8g.cluster"
  node_num = 4
  zone_list = data.ctyun_zones.test.zones
  disk_size = 100
  disk_type = "SAS"
  vpc_id= local.real_vpc_id
  subnet_id = local.real_subnet_id
  security_group_id= local.real_security_group_id
  cycle_type = "month"
  cycle_count = 1
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