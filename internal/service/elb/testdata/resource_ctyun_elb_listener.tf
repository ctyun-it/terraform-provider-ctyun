resource "ctyun_elb_listener" "%[1]s"{
  loadbalancer_id = "%[2]s"
  name = "%[3]s"
  protocol = "%[4]s"
  protocol_port = %[5]d
  default_action_type = "%[6]s"
  %[7]s
  %[8]s
  %[9]s
  %[10]s
  %[11]s
  %[12]s
  %[13]s
  status = "%[14]s"
}

