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

# 创建 VPC
resource "ctyun_vpc" "example" {
  name        = "tf-vpc-for-rocketmq"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  region_id   = "200000001703"
}

# 创建子网
resource "ctyun_subnet" "example" {
  name        = "tf-subnet-for-rocketmq"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  vpc_id      = ctyun_vpc.example.id
  region_id   = "200000001703"
  dns         = ["8.8.8.8", "8.8.4.4"]
}

# 创建安全组
resource "ctyun_security_group" "example" {
  name        = "tf-sg-for-rocketmq"
  description = "terraform测试使用"
  vpc_id      = ctyun_vpc.example.id
  region_id   = "200000001703"
}

# 查询可用区
data "ctyun_zones" "test" {
  region_id = "200000001703"
}

# RocketMQ 实例
resource "ctyun_rocketmq_instance" "example" {
  instance_name     = "tf-rocketmq-instance"
  spec_name         = "rocketmq.4u8g.cluster"
  node_num          = 4
  zone_list         = data.ctyun_zones.test.zones
  disk_size         = 100
  disk_type         = "SAS"
  vpc_id            = ctyun_vpc.example.id
  subnet_id         = ctyun_subnet.example.id
  security_group_id = ctyun_security_group.example.id
  cycle_type        = "month"
  cycle_count       = 1
  region_id         = "200000001703"
}

# 查询 RocketMQ 实例下的所有 Topic
data "ctyun_rocketmq_topics" "test" {
  instance_id = ctyun_rocketmq_instance.example.id
}
