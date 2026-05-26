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

# 从 ZOS 桶上传代码包创建函数
resource "ctyun_function" "function_from_zos" {
  name                          = "tf-function-from-zos"
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
  description                   = "Terraform 创建的测试函数 - 从 ZOS 上传"
  code_bucket                   = "your-bucket-name"
  code_key                      = "your-function-code.zip"
  
  # 可选配置
  # environment = {
  #   ENVIRONMENT = "production"
  #   LOG_LEVEL   = "info"
  # }
  # role = "your-iam-role-name"
}

# 使用内联代码创建函数
resource "ctyun_function" "function_inline_code" {
  name                          = "tf-function-inline"
  create_type                   = 1
  runtime_runtime               = "Python3.9"
  runtime_handle_type           = "http"
  runtime_handler               = "main.handler"
  runtime_execute_timeout       = 30
  runtime_instance_concurrency  = 5
  container_time_zone           = "UTC"
  container_disk_size           = 512
  container_memory_size         = 128
  container_cpu                 = 0.1
  container_listen_port         = 8080
  description                   = "Terraform 创建的测试函数 - 内联代码"
  code_content                  = base64encode(file("${path.module}/function_code.zip"))
  
  # 自定义容器配置
  # container_max_scale = 100
  # container_fast_start = 1
  
  # 环境变量配置
  environment = {
    TF_MANAGED    = "true"
    CREATED_BY    = "terraform"
    PROJECT_NAME  = "faas-demo"
  }
}

