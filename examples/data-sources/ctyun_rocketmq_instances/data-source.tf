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

# 查询所有 RocketMQ 实例
data "ctyun_rocketmq_instances" "test" {
}
