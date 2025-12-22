resource "ctyun_postgresql_backup" "%[1]s" {
  project_id = "%[2]s"
  instance_id = "%[3]s"
  name = "%[4]s"
  description = "%[5]s"
}
