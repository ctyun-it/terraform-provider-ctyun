resource "ctyun_ecs_data_volume" "%[1]s" {
  instance_id = "%[2]s"
  ebs_ids = ["%[3]s"]
}