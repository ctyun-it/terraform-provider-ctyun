resource "ctyun_nat_snat" "%[1]s"{
    nat_gateway_id = "%[2]s"
    source_cidr = "%[3]s"
    snat_ips = "%[4]s"
}
