resource "ctyun_nat_dnat" "%[1]s"{
  nat_gateway_id = "%[2]s"
  external_id = "%[3]s"
  protocol = "%[4]s"
  external_port = %[5]d
  internal_port = %[6]d
  dnat_type = "%[7]s"
  server_type = "%[8]s"
  instance_id = "%[9]s"
}
