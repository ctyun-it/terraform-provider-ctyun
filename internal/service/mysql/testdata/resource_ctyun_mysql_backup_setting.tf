resource "ctyun_mysql_backup_setting" "%[1]s" {
  instance_id                    = "%[2]s"
  storage_day                = %[3]d
  allow_earliest_time        = "%[6]s"
  trigger_days_of_week       = %[7]s
}
