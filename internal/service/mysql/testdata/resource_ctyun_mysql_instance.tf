resource "ctyun_mysql_instance" "%[1]s" {
  cycle_type = "%[2]s"
  vpc_id = "%[3]s"
  subnet_id = "%[4]s"
  security_group_id = "%[5]s"
  name = "%[6]s"
  password = "%[7]s"
  %[8]s // cycle_count
  %[9]s // auto_renew
  flavor_name = "%[10]s"
  prod_id = "%[11]s"
  %[12]s  //write_port
  storage_type = "%[13]s"
  storage_space = %[14]d
  %[15]s  // availability_zone_info
  %[16]s // running_control
  %[17]s // backup_storage_space
}

