resource "ctyun_ec_cda_instance" "%[1]s" {
  ec_id             = "%[2]s"
  cgw_id            = "%[3]s"
  cda_id            = "%[4]s"
  cda_name          = "%[5]s"
  cda_cidr_v4_list  = [%[6]s]
  rtb_id            = "%[7]s"
  cda_info          = %[8]s
  account           = "%[9]s"
  weights           = %[10]d
  route_learn       = %[11]d
  route_sync        = %[12]d
}