data "ctyun_rabbitmq_exchanges" "%[1]s" {
  instance_id = "%[2]s"
  vhost = "/"
  name = "%[3]s"
}