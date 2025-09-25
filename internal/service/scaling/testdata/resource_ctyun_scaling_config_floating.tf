# 使用eip方案
resource "ctyun_scaling_config" "%[1]s" {
  name            = "%[2]s"
  image_id        = "%[3]s"
  flavor_name     = "%[4]s"
  use_floatings   = "%[5]s"
  bandwidth       = %[6]d
  login_mode      = "%[7]s"
  password        = "%[8]s"
  monitor_service = %[9]t
  az_names        = %[10]s
  tags            = %[11]s
  volumes         = %[12]s
}

