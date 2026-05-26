resource "ctyun_function" "%[1]s" {
  name                          = "%[2]s"
  create_type                   = 1
  runtime_runtime               = "%[3]s"
  runtime_handle_type           = "%[4]s"
  runtime_handler               = "%[5]s"
  runtime_execute_timeout       = %[6]d
  runtime_instance_concurrency  = %[7]d
  container_time_zone           = "%[8]s"
  container_disk_size           = %[9]d
  container_memory_size         = %[10]d
  container_cpu                 = %[11]f
  container_listen_port         = %[12]d
  description                   = "%[13]s"
  code_bucket                   = "%[14]s"
  code_key                      = "%[15]s"
}


