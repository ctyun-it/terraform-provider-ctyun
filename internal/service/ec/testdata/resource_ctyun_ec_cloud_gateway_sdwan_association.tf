resource "ctyun_ec_cloud_gateway_sdwan_association" "%[1]s" {
  ec_id    = "%[2]s"
  sdwan_id = "%[3]s"
  cgw_list = [%[4]s]
}