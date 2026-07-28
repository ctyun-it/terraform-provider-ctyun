resource "ctyun_postgresql_white_list" "%[1]s" {
  instance_id    = "%[2]s"
  mode       = "%[3]s"
  ip_list    = %[4]s
}