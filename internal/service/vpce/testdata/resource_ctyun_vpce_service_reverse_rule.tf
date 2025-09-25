resource "ctyun_vpce_service_reverse_rule" "%[1]s" {
  endpoint_service_id = "%[2]s"
  endpoint_id = "%[3]s"
  transit_ip = "%[4]s"
  transit_port = 1
  target_ip = "%[5]s"
  target_port = 2
  protocol = "TCP"
}
