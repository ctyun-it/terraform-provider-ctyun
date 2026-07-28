resource "ctyun_function_trigger" "%[1]s" {
  function_name = "%[2]s"
  trigger_name  = "%[3]s"
  trigger_type  = "rocketmq"
  region_id     = "%[4]s"

  event_data = jsonencode({
    instanceId = "%[5]s"
    urlPath    = ""
    topic      = "test"
    tag        = ""
    consumerGroup = {
      autoCreate = true
      groupId    = "group"
    }
    consumeFrom       = 0
    consumeTimestamp  = ""
    synchronize       = true
    concurrencyLimit  = 0
    enable            = true
    auth = {
      accessKey     = "dddd"
      secretKey     = "Fep7GFky9FYW+sAaDQi3ZA=="
      securityToken = ""
    }
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
