resource "ctyun_mysql_backup_setting" "%[1]s" {
  instance_id                    = "%[2]s"
  storage_day                = %[3]d
  frequency_backup           = %[4]t
  frequency_backup_unit_time = %[5]d
  allow_earliest_time        = "%[6]s"
  trigger_days_of_week       = %[7]s
}
