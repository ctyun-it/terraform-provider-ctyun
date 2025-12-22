data "ctyun_postgresql_backups" "%[1]s" {
  instance_id = "%[2]s"
  name        = "%[3]s"
  type        = "%[4]s"
}
