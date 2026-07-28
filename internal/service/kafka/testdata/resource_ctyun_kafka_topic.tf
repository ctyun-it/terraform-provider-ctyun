resource "ctyun_kafka_topic" "%[1]s" {
  name = "%[2]s"
  instance_id = "%[3]s"
  partition_num  = %[4]d
}
