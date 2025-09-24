resource "ctyun_rabbitmq_instance" "%[1]s" {
  instance_name = "%[2]s"
  spec_name = "%[3]s"
  node_num = %[4]d
  zone_list = ["%[5]s"]
  disk_size = %[6]d
  disk_type = "%[7]s"
  vpc_id = "%[8]s"
  subnet_id = "%[9]s"
  security_group_id = "%[10]s"
  cycle_type = "month"
  cycle_count = 1
}
