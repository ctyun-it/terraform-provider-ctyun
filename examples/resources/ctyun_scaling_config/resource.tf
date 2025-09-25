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

variable "password" {
  type      = string
  sensitive = true
}

data "ctyun_images" "image_test" {
  name       = "CentOS Linux 8.4"
  visibility = "public"
  page_no = 1
  page_size = 10
}

locals {
  image_id = data.ctyun_images.image_test.images[0].id
}

resource "ctyun_scaling_config" "scaling_config_example" {
  name            = "sc-example"
  image_id        = local.image_id
  flavor_name     = "s7.large.2"
  use_floatings   = "auto"
  bandwidth       = 1
  login_mode      = "password"
  password        = var.password
  monitor_service = true
  az_names        = ["cn-huadong1-jsnj1A-public-ctcloud", "cn-huadong1-jsnj2A-public-ctcloud"]
  tags            =[{"key":"provider", "value":"scaling_conifg"}, {"key":"version", "value":"1.1.1"}]
  volumes         = [{"volume_type":"SATA", "volume_size":40, "flag":"OS"}, {"volume_type":"SAS", "volume_size":100, "flag":"DATA"}]
}