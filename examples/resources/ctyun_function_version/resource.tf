terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 可参考 index.md，在环境变量中配置 ak、sk、资源池 ID、可用区名称
provider "ctyun" {
  env = "prod"
}

# 首先创建一个函数
resource "ctyun_function" "test_function" {
  name                          = "tf-function-for-version"
  create_type                   = 1
  runtime_runtime               = "Python3.9"
  runtime_handle_type           = "event"
  runtime_handler               = "index.handler"
  runtime_execute_timeout       = 60
  runtime_instance_concurrency  = 10
  container_time_zone           = "Asia/Shanghai"
  container_disk_size           = 512
  container_memory_size         = 256
  container_cpu                 = 0.25
  container_listen_port         = 8080
  description                   = "用于版本测试的函数"
  code_content                  = base64encode(file("${path.module}/function_code.zip"))
}

# 创建函数版本
resource "ctyun_function_version" "version_v1" {
  function_name = ctyun_function.test_function.name
  description   = "v1.0 - 初始版本"
  
  # region_id 可选，默认从 provider 或环境变量获取
  # region_id = "200000002401"
}

# 创建第二个版本
resource "ctyun_function_version" "version_v2" {
  function_name = ctyun_function.test_function.name
  description   = "v2.0 - 功能更新版本"
}

