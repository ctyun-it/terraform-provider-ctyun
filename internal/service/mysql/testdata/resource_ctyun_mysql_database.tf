resource "ctyun_mysql_database" "%[1]s" {
  instance_id      = "%[2]s"
  name         = "%[3]s"
  charset_name = %[4]s
  description  = "%[5]s"
}
