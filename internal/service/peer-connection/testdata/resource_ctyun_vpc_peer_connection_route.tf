resource "ctyun_vpc_route_table_rule" "%[1]s" {
  ip_version     = %[2]d
  next_hop_id    = "%[3]s"
  destination    = "%[4]s"
  next_hop_type  = "vpcpeering"
  route_table_id = "%[5]s"
  description    = "%[6]s"
}

