resource "ctyun_vpce" "%[1]s" {
  name  = "%[2]s"
  endpoint_service_id = "%[7]s"
  vpc_id = "%[5]s"
  subnet_id = "%[6]s"
  whitelist_flag = "%[3]s"
  %[4]s
}
