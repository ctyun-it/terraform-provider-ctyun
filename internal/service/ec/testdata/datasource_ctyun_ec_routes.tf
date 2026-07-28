data "ctyun_ec_routes" "%[1]s" {
  ec_id  = "%[2]s"
  cgw_id = "%[3]s"
  rtb_id = "%[4]s"
}