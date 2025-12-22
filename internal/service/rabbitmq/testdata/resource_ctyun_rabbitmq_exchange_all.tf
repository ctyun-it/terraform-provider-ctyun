resource "ctyun_rabbitmq_exchange" "%[1]s" {
  instance_id = "%[2]s"
  vhost = "/"
  name = "%[3]s"
  type = "%[4]s"
  x_delayed_type = "%[5]s"
  durable = true
  auto_delete = true
  internal = true
  alternate_exchange = "%[6]s"
}
