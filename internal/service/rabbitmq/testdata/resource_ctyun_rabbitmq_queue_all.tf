resource "ctyun_rabbitmq_queue" "%[1]s" {
  instance_id = "%[2]s"
  vhost = "/"
  name = "%[3]s"
  durable = true
  auto_delete = true
  x_queue_mode = "%[4]s"
  x_overflow = "%[5]s"
  x_dead_letter_exchange = "%[6]s"
  x_dead_letter_routing_key = "%[7]s"
  x_message_ttl = %[8]d
  x_max_length = %[9]d
  x_expires = %[10]d
  x_max_priority = %[11]d
}


