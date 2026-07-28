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

resource "ctyun_rabbitmq_queue" "test" {
  instance_id  = "8ccc8af2e6704080a72548735a081660"
  vhost        = "/"
  name         = "example-queue"
  durable      = true
  auto_delete  = true
  x_queue_mode = "default"
  x_expires    = 86400000
  # x_dead_letter_exchange = "dead-letter-exchange"
  x_dead_letter_routing_key = "dead-letter-routing-key"
  x_message_ttl             = 3600000
  x_max_length              = 1000
  x_overflow                = "drop-head"
  x_max_priority            = 5
}

