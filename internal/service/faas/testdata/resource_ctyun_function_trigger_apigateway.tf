resource "ctyun_function_trigger" "%[1]s" {
  function_name = "%[2]s"
  trigger_name  = "trigger-apigateway-%[3]s"
  trigger_type  = "apigateway"
  region_id     = "%[4]s"

  event_data = jsonencode({
    gatewayInstanceId = "%[5]s"
    vpceId            = "%[6]s"
    domain            = ["example.com"]
    path              = "/api/v1/*"
    methods           = ["GET", "POST"]
    priority          = 10
  })

  enable = true
}
