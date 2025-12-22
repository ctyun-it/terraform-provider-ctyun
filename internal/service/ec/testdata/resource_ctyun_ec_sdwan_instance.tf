resource "ctyun_ec_sdwan_instance" "%[1]s" {
  ec_id       = "%[2]s"
  cgw_id      = "%[3]s"
  sdwan_id    = "%[4]s"
  rtb_id      = "%[5]s"
  weights     = %[6]d
  route_learn = %[7]d
  route_sync  = %[8]d
}
