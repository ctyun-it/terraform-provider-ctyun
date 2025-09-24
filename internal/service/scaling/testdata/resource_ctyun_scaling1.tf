resource "ctyun_scaling_group" "%[1]s" {
  security_group_id_list = %[2]s
  name                   = "%[3]s"
  health_mode            = "%[4]s"
  subnet_id_list         = %[5]s
  move_out_strategy      = "%[6]s"
  vpc_id                 = "%[7]s"
  min_count              = %[8]d
  max_count              = %[9]d
  health_period          = %[10]d
  use_lb                 = %[11]d
  lb_list                = %[12]s
  config_list            = %[13]s
  az_strategy            = "%[14]s"
  status                 = "%[15]s"
  delete_protection      = "%[16]s"
}

