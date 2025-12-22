resource "ctyun_rabbitmq_queue" "%[1]s" {
  instance_id = "%[2]s"
  vhost = "/"
  name = "%[3]s"
}
