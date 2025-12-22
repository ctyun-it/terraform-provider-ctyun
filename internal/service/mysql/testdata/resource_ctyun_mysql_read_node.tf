resource "ctyun_mysql_readonly_instance" "%[1]s" {
  instance_id      = "%[2]s"
  cycle_type   = "%[3]s"
  flavor_name  = "%[4]s"
  project_id   = "%[5]s"
  storage_type = "%[6]s"
  storage_space = %[7]d
  name         = "%[8]s"
}