resource "ctyun_mysql_instance" "%[1]s" {
  cycle_type        = "%[2]s"
  vpc_id            = "%[3]s"
  subnet_id         = "%[4]s"
  security_group_id = "%[5]s"
  name              = "%[6]s"
  password          = "%[7]s"
  %[8]s
  %[9]s
  flavor_name       = "%[10]s"
  prod_id           = "%[11]s"
  write_port        = %[12]d
  storage_type      = "%[13]s"
  storage_space     = %[14]d
  availability_zone_info = %[15]s
  backup_storage_space = %[16]d
  %[17]s // running_control
}
