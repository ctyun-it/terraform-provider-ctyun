terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 可参考index.md，在环境变量中配置ak、sk、资源池ID、可用区名称
provider "ctyun" {
  env = "prod"
}

resource "ctyun_elb_health_check" "health_check_test" {
  name     = "tf_health_check"
  protocol = "HTTP"
  timeout = 60
  interval = 60
  max_retry = 10
  http_method = "POST"
  http_url_path = "/health"
  http_expected_codes = ["http_2xx","http_3xx","http_4xx","http_5xx"]
  protocol_port = 8080
}
