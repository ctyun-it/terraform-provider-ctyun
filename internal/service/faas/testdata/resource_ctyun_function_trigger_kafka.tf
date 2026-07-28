resource "ctyun_function_trigger" "%[1]s" {
  function_name = "%[2]s"
  trigger_name  = "%[3]s"
  trigger_type  = "kafka"
  region_id     = "%[4]s"

  event_data = jsonencode({
    instanceId = "%[5]s"
    urlPath    = ""
    topic      = "test"
    consumerGroup = {
      autoCreate = true
      groupId    = ""
    }
    consumerNum      = 1
    offset           = 1
    synchronize      = true
    concurrencyLimit = 1
    enable           = true
    pushConfig = {
      batchMaxSize = 1
      interval     = 0
      format       = 0
    }
    strategyConfig = {
      retryStrategy        = "exponential"
      faultToleranceStrategy = "on"
      deadletterStrategy = {
        enable = false
        instanceType = ""
        instanceId   = ""
        topic        = ""
        auth = {
          accessKey = ""
          secretKey = ""
        }
      }
    }
  })

  enable = true
}
