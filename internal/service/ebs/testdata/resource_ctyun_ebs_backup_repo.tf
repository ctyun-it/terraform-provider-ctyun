resource "ctyun_ebs_backup_repo" "%[1]s" {
  name = "%[2]s"
  size = %[3]d
  cycle_count = "5"
  cycle_type  = "MONTH"
}
