resource "ctyun_ec_vpc_instance" "%[1]s" {
  ec_id       = "%[2]s"
  cgw_id      = "%[3]s"
  rtb_id      = "%[4]s"
  vpc_id      = "%[5]s"
  route_learn = %[6]d
  route_sync  = %[7]d
  subnets     = [%[8]s]
}
