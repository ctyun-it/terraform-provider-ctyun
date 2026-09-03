resource "ctyun_nat" "%[1]s"{
    vpc_id = "%[2]s"
    spec = "%[3]s"
    name = "%[4]s"
    description = "%[5]s"
    cycle_type = "%[6]s"
    %[7]s
    tcp_expire_time = %[8]d
    udp_expire_time = %[9]d
    icmp_expire_time = %[10]d
    tcp_delay_close_time = %[11]d
}
