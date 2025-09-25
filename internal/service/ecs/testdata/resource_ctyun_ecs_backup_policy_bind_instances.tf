resource "ctyun_ecs_backup_policy_bind_instances" "%[1]s" {
  policy_id = %[2]s
  instance_id_list = "%[3]s"
}
