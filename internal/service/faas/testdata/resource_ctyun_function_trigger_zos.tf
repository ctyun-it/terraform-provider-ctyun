resource "ctyun_function_trigger" "%[1]s" {
  function_name = "%[2]s"
  trigger_name  = "trigger-zos-%[3]s"
  trigger_type  = "zos"
  region_id     = "%[4]s"

  event_data = jsonencode({
    type          = ["s3:ObjectCreated:Put"]
    bucket        = "%[5]s"
    objectPrefix  = ["prefix1"]
    objectSuffix  = [".jpg", ".png"]
    syncronize    = true
  })

  enable = true
}
