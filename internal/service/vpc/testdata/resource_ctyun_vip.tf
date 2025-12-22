resource "ctyun_vip" "%[1]s" {
  subnet_id  = "%[2]s"
  vpc_id     = "%[3]s"
  ip_address = "%[4]s"
  vip_type   = "%[5]s"

}