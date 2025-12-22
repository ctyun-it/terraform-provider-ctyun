resource "ctyun_mysql_backup" "%[1]s" {
  instance_id     = "%[2]s"
  project_id  = "%[3]s"
  description = "%[4]s"
  task_type   = "%[5]s"
}
