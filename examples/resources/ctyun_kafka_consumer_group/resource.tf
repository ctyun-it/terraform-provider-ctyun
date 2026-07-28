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

resource "ctyun_kafka_consumer_group" "tbidgqvfbs" {
  name        = "test_kafka_consumer_group"
  instance_id = "4bd607df61d348b1949db223614315c1"
  description = "desc"
}

