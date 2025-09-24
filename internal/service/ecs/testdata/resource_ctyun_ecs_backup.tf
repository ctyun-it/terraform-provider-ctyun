resource "ctyun_ecs_backup" "%[1]s" {
  repository_id = "%[2]s"
  instance_id = "%[3]s"
  name  = "%[4]s"
  full_backup = false
}
