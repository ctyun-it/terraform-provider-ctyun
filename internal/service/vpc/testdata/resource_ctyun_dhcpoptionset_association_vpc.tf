resource "ctyun_dhcpoptionset_association_vpc" "%[1]s" {
  dhcp_option_sets_id = "%[2]s"
  vpc_ids             = [%[3]s]
}
