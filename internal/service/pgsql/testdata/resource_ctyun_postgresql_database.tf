resource "ctyun_postgresql_database" "%[1]s" {
  instance_id         = "%[2]s"
  name            = "%[3]s"
  charset_name    = "%[4]s"
  charset_collate = "%[5]s"
  charset_type    = "%[6]s"
  owner           = "%[7]s"
  description     = "%[8]s"
}
