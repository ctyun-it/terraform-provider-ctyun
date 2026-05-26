resource "ctyun_function_trigger" "%[1]s" {
  function_name = "%[2]s"
  trigger_name  = "%[3]s"
  trigger_type  = "schedule"
  event_data = jsonencode({
    cronExpr = "0 */5 * * * ?"
    data     = "{\"key\":\"value\"}"
  })

  enable = true
}
