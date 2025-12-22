resource "ctyun_mysql_database" "%[1]s" {
  instance_id      = "%[2]s"
  project_id   = "%[3]s"
  name         = "%[4]s"
  charset_name = %[5]s
  description  = "%[6]s"
}
