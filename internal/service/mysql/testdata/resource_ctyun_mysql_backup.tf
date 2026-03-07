resource "ctyun_mysql_backup" "%[1]s" {
  instance_id     = "%[2]s"
  description = "%[3]s"
  task_type   = "%[4]s"
}
