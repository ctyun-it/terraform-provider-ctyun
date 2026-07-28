resource "ctyun_crs_vpc_attach" "%[1]s" {
  vpc_id = "%[2]s"
  subnet_id = "%[3]s"
}