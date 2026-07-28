resource "ctyun_ec_region_peer" "%[1]s" {
  name        = "%[2]s"
  ec_id       = "%[3]s"
  src_cgw_id  = "%[4]s"
  dst_cgw_id  = "%[5]s"
  packet_id   = "%[6]s"
  rate        = %[7]d
  route_learn = %[8]d
}

