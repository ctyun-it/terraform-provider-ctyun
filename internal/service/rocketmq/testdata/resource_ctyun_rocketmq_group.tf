resource "ctyun_rocketmq_group" "%[1]s" {
  instance_id                  = "%[2]s"
  name                         = "%[3]s"
  consume_enable               = %[4]t
  first_consume_mechanism      = %[5]d
  pull_mechanism               = %[6]d
  retry_max_times              = %[7]d
  remark                       = "%[8]s"
  broker_names                 = ["broker_1"]
}
