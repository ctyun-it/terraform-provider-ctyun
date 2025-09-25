resource "ctyun_elb_target" "%[1]s" {
  target_group_id = "%[2]s"
  instance_type = "%[3]s"
  instance_id = "%[4]s"
  protocol_port = %[5]d
  %[6]s
}
