resource "ctyun_ec_route" "%[1]s" {
  ec_id               = "%[2]s"
  cgw_id             = "%[3]s"
  rtb_id              = "%[4]s"
  cidr                = "%[5]s"
  ip_version          = "%[6]s"
  description         = "%[7]s"
  is_black_hole_route = %[8]t
  next_hop_type       = "%[9]s"
  next_hop_id         = "%[10]s"
}

