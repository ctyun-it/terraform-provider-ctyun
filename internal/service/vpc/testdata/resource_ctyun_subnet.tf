resource "ctyun_subnet" "%[1]s" {
  vpc_id = "%[5]s"
  name        = "%[2]s"
  cidr        = "192.168.1.0/24"
  description = "%[3]s"
  dns         = [
    "%[4]s",
  ]
  enable_ipv6 = true
}
