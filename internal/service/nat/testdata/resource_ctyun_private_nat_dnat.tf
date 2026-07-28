


resource "ctyun_private_nat_dnat" "%[1]s"{
    nat_gateway_id = "%[2]s"
    external_ip = "%[3]s"
    protocol = "%[4]s"
    external_port = %[5]d
    internal_port = %[6]d
    internal_ip = "%[7]s"
}
