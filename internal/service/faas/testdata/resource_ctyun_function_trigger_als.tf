resource "ctyun_function_trigger" "%[1]s" {
  function_name = "%[2]s"
  trigger_name  = "trigger-als-%[3]s"
  trigger_type  = "als"
  region_id     = "%[4]s"

  event_data = jsonencode({
    logProjectCode = "%[5]s"
    logProjectName = "%[6]s"
    logUnitCode    = "%[7]s"
    logUnittName   = "%[8]s"
    syncronize     = true
  })

  enable = true
}
