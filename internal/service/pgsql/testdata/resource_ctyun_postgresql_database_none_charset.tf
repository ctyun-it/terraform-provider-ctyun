resource "ctyun_postgresql_database" "%[1]s" {
  project_id   = "%[2]s"
  instance_id      = "%[3]s"
  name         = "%[4]s"
  charset_name = "%[5]s"
  owner        = "%[6]s"
  description  = "%[7]s"
}
