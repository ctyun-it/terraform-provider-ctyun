resource "ctyun_nat" "%[1]s"{
    vpc_id = "%[2]s"
    spec = "%[3]s"
    name = "%[4]s"
    description = "%[5]s"
    cycle_type = "%[6]s"
    %[7]s
}
