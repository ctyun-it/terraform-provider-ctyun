resource "ctyun_function_trigger" "%[1]s" {
  function_name = "%[2]s"
  trigger_name  = "trigger-rabbitmq-%[3]s"
  trigger_type  = "rabbitmq"
  region_id     = "%[4]s"

  event_data = jsonencode({
    instanceId  = "%[5]s"
    vhost       = "/"
    queueName   = "test-queue"
    routingKey  = "test.key"
    synchronize = true
  })

  enable = true
}
