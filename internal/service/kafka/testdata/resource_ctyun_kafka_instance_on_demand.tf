resource "ctyun_kafka_instance" "%[1]s" {
  instance_name = "%[2]s"
  engine_version = "%[3]s"
  spec_name = "%[4]s"
  node_num = %[5]d
  zone_list = ["%[6]s"]
  disk_type = "%[7]s"
  disk_size = %[8]d
  vpc_id = "%[9]s"
  subnet_id = "%[10]s"
  security_group_id = "%[11]s"
  retention_hours = %[12]d
  cycle_type = "on_demand"
  plain_port  = 9001
  sasl_port = 9002
  ssl_port = 9003
  http_port = 9004
}
