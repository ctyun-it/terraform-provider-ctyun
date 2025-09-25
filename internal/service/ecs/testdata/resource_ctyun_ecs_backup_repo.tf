resource "ctyun_ecs_backup_repo" "%[1]s" {
  name = "%[2]s"
  cycle_count = "5"
  cycle_type  = "MONTH"
}
