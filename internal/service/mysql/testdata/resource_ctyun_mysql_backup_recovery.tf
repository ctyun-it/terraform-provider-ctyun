resource "ctyun_mysql_backup_recovery" "%[1]s" {
  instance_id      = "%[2]s"
  project_id   = "%[3]s"
  src_instance_id  = "%[4]s"
  dst_instance_id  = "%[5]s"
  to_timepoint = "%[6]s"
}
