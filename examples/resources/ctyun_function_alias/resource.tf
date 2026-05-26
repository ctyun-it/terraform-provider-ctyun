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
  name                          = "tf-function-for-alias"
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
  description                   = "用于别名测试的函数"
  code_content                  = base64encode(file("${path.module}/function_code.zip"))
}

# 创建函数版本
resource "ctyun_function_version" "test_version" {
  function_name = ctyun_function.test_function.name
  description   = "v1.0 - 用于别名测试"
}

# 创建函数别名（简单配置）
resource "ctyun_function_alias" "alias_prod" {
  function_name = ctyun_function.test_function.name
  alias_name    = "prod"
  version_id    = ctyun_function_version.test_version.id
  description   = "生产环境别名"
  
  # region_id 可选，默认从 provider 或环境变量获取
  # region_id = "200000002401"
}

# 创建带灰度发布的别名
resource "ctyun_function_alias" "alias_with_gray" {
  function_name   = ctyun_function.test_function.name
  alias_name      = "gray-release"
  version_id      = ctyun_function_version.test_version.id
  description     = "带灰度发布的别名"
  
  # 灰度发布配置
  gray_type       = 1          # 按百分比随机灰度
  gray_version_id = ctyun_function_version.test_version.id
  gray_weight     = 10         # 10% 的流量到灰度版本
  
  depends_on = [ctyun_function_version.test_version]
}
