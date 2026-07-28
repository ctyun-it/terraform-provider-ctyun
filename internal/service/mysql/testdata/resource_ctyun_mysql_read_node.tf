resource "ctyun_mysql_readonly_instance" "%[1]s" {
  instance_id      = "%[2]s"
  cycle_type   = "%[3]s"
  flavor_name  = "%[4]s"
  storage_type = "%[5]s"
  storage_space = %[6]d
  name         = "%[7]s"
}