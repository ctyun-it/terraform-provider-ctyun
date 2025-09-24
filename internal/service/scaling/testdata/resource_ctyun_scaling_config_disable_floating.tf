# 不使用eip方案
resource "ctyun_scaling_config" "%[1]s" {
  name            = "%[2]s"
  image_id        = "%[3]s"
  flavor_name     = "%[4]s"
  use_floatings   = "%[5]s"
  login_mode      = "%[6]s"
  key_pair_id     = "%[7]s"
  monitor_service = %[8]t
  az_names        = %[9]s
  tags            = %[10]s
  volumes         = %[11]s
}

