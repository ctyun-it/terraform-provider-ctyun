resource "ctyun_vpc" "%[1]s" {
  name        = "%[2]s"
  description = "%[3]s"
  cidr        = "%[4]s"
  enable_ipv6 = true
}
