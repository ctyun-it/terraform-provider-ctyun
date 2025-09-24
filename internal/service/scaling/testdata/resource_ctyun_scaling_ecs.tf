resource "ctyun_scaling_ecs" "%[1]s" {
  group_id           = %[2]d
  instance_uuid_list = %[3]s
  protect_status     = "%[4]s"
}