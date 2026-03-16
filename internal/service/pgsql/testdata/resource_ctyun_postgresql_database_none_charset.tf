resource "ctyun_postgresql_database" "%[1]s" {
  instance_id      = "%[2]s"
  name         = "%[3]s"
  charset_name = "%[4]s"
  owner        = "%[5]s"
  description  = "%[6]s"
}
