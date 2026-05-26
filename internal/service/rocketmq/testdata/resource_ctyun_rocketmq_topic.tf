resource "ctyun_rocketmq_topic" "%[1]s" {
  instance_id      = "%[2]s"
  name             = "%[3]s"
  queue_nums = 8
  broker_names     = [ "broker_1"]
  order            = %[4]t
  perm             = %[5]d
  remark           = "%[6]s"
}
