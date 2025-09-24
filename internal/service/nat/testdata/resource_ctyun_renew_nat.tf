resource "ctyun_nat" "%[1]s"{
    vpc_id = "%[2]s"
    spec = "%[3]s"
    name = "%[4]s"
    description = "%[5]s"
    cycle_type = "%[6]s"
    cycle_count = "%[7]s"
    az_name = "%[8]s"
}
