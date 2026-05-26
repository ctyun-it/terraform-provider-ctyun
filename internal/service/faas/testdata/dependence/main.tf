resource "ctyun_function" "hello_server" {
  name                         = "faas_hello_server"
  create_type                  = 1
  runtime_runtime              = "python3.9"
  runtime_handle_type          = "event"
  runtime_handler              = "index.handler"
  runtime_execute_timeout      = 50
  runtime_instance_concurrency = 50
  container_time_zone          = "UTC"
  container_disk_size          = 512
  container_memory_size        = 256
  container_cpu                = 0.2
  container_listen_port        = 8080
  description                  = "Hello Server Function"
  code_bucket                  = "bucket-for-faas"
  code_key                     = "hello_server.zip"
}
