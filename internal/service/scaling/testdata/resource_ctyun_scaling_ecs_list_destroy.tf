resource "ctyun_scaling_group" "%[1]s" {
  security_group_id_list = %[2]s
  name                   = "%[3]s"
  health_mode            = "%[4]s"
  subnet_id_list         = %[5]s
  move_out_strategy      = "%[6]s"
  vpc_id                 = "%[7]s"
  min_count              = %[8]d
  max_count              = %[9]d
  expected_count         = %[10]d
  health_period          = %[11]d
  use_lb                 = %[12]d
  config_list            = %[13]s
  add_instance_uuid_list     = %[14]s
  remove_instance_uuid_list  = %[15]s
  az_strategy            = "%[16]s"
}

