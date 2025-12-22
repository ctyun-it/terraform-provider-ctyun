resource "ctyun_rabbitmq_exchange" "%[1]s" {
  instance_id = "%[2]s"
  vhost = "/"
  name = "%[3]s"
  type = "%[4]s"
}
