resource "ctyun_postgresql_white_list" "%[1]s" {
  project_id = "%[2]s"
  instance_id    = "%[3]s"
  mode       = "%[4]s"
  ip_list    = %[5]s
}