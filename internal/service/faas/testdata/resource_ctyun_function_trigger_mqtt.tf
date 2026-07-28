resource "ctyun_function_trigger" "%[1]s" {
  function_name = "%[2]s"
  trigger_name  = "trigger-mqtt-%[3]s"
  trigger_type  = "mqtt"
  region_id     = "%[4]s"

  event_data = jsonencode({
    instanceId  = "%[5]s"
    topic       = "test/topic"
    synchronize = true
  })

  enable = true
}
