resource "ctyun_mysql_backup_recovery" "%[1]s" {
  instance_id      = "%[2]s"
  src_instance_id  = "%[3]s"
  dst_instance_id  = "%[4]s"
  to_timepoint = "%[5]s"
}
