data "ctyun_rabbitmq_instances" "%[1]s" {
  instance_id = %[2]s
}

data "ctyun_rabbitmq_instances" "test" {
  page_no = 1
  page_size = 10
}