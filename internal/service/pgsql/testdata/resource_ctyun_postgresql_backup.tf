resource "ctyun_postgresql_backup" "%[1]s" {
  instance_id = "%[2]s"
  name = "%[3]s"
  description = "%[4]s"
}
