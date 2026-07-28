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
  name                          = "tf-function-for-domain"
  create_type                   = 1
  runtime_runtime               = "Python3.9"
  runtime_handle_type           = "http"
  runtime_handler               = "main.handler"
  runtime_execute_timeout       = 60
  runtime_instance_concurrency  = 10
  container_time_zone           = "Asia/Shanghai"
  container_disk_size           = 512
  container_memory_size         = 256
  container_cpu                 = 0.25
  container_listen_port         = 8080
  description                   = "用于自定义域名测试的函数"
  code_content                  = base64encode(file("${path.module}/function_code.zip"))
}

# 创建函数版本
resource "ctyun_function_version" "test_version" {
  function_name = ctyun_function.test_function.name
  description   = "v1.0 - 用于域名测试"
}

# 创建自定义域名（HTTP 协议）
resource "ctyun_function_domain" "domain_http" {
  domain_name   = "api.example.com"
  protocol      = "HTTP"
  description   = "HTTP 协议的自定义域名"
  cname_check   = false
  
  # 路由配置
  route_config {
    routes {
      function_name = ctyun_function.test_function.name
      path          = "/api/*"
      qualifier     = ctyun_function_version.test_version.id
      
      # HTTP 方法配置
      methods = ["GET", "POST", "PUT", "DELETE"]
      
      # JWT 认证配置
      enable_jwt = 0
    }
    
    routes {
      function_name = ctyun_function.test_function.name
      path          = "/health"
    }
  }
  
  # region_id 可选，默认从 provider 或环境变量获取
  # region_id = "200000002401"
}

# 创建自定义域名（HTTPS 协议，带证书配置）
resource "ctyun_function_domain" "domain_https" {
  domain_name   = "secure-api.example.com"
  protocol      = "HTTPS"
  description   = "HTTPS 协议的自定义域名，带 SSL 证书"
  cname_check   = true
  
  # SSL 证书配置
  cert_config {
    cert_name    = "example-com-cert"
    certificate  = file("${path.module}/ssl_certificate.pem")
    private_key  = file("${path.module}/ssl_private_key.pem")
  }
  
  # 路由配置
  route_config {
    routes {
      function_name = ctyun_function.test_function.name
      path          = "/api/v1/*"
      qualifier     = ctyun_function_version.test_version.id
      methods       = ["GET", "POST"]
      enable_jwt    = 1
    }
  }
  
  # JWT 认证配置
  auth_config {
    auth_type = "JWT"
    
    jwt_config {
      jwks = file("${path.module}/jwks.json")
      
      token_config {
        location = "HEADER"
        name     = "Authorization"
      }
      
      match_mode {
        mode = "All"
        path = ["/api/*"]
      }
    }
  }
  
  depends_on = [ctyun_function_version.test_version]
}

