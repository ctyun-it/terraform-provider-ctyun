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

# 查询 RocketMQ 可用规格
data "ctyun_rocketmq_specs" "test" {
}
