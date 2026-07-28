resource "ctyun_dhcpoptionset" "%[1]s" {
  name         = "%[1]s"
  description  = "%[2]s"
  domain_name  = "%[3]s"
  dns_list     = [%[4]s]
}
