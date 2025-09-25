resource "ctyun_nat_snat" "%[1]s"{
    nat_gateway_id = "%[2]s"
    %[3]s
    snat_ips = %[4]s
    description = "%[5]s"
}
