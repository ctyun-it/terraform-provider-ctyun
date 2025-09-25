resource "ctyun_scaling_ecs_protection" "%[1]s" {
  group_id         = %[2]s
  instance_id_list = %[3]s
  protect_status   = %[4]t
}