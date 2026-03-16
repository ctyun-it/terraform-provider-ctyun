resource "ctyun_postgresql_database" "%[1]s" {
  project_id      = "%[2]s"
  instance_id         = "%[3]s"
  name            = "%[4]s"
  charset_name    = "%[5]s"
  charset_collate = "%[6]s"
  charset_type    = "%[7]s"
  owner           = "%[8]s"
  description     = "%[9]s"
}
