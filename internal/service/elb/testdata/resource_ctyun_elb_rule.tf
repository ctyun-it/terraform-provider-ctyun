resource "ctyun_elb_rule" "%[1]s" {
  listener_id = "%[2]s"
  conditions  = [%[3]s]
  action_type = "%[4]s"
  action_target_groups = [%[5]s]
}

