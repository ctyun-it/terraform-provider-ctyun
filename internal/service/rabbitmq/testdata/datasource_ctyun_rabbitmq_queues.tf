data "ctyun_rabbitmq_queues" "%[1]s" {
  instance_id = "%[2]s"
  vhost = "/"
  name = "%[3]s"
}