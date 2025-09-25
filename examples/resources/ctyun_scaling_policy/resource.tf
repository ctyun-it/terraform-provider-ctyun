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


variable "key_pair" {
  type      = string
  sensitive = true
}

resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-scaling"
  cidr        = "192.168.0.0/16"
  description = "terraform-kafka测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-scaling"
  cidr        = "192.168.1.0/24"
  description = "terraform-scaling测试使用"
  dns = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
}
resource "ctyun_security_group" "sg_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-scaling"
  description = "terraform-scaling测试使用"
  lifecycle {
    prevent_destroy = false
  }
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

resource "ctyun_keypair" "scaling_test" {
  name       = "key-pair-scaling-test"
  public_key = var.key_pair
}

resource "ctyun_scaling_config" "config_test" {
  name            = "sc-for-policy"
  image_id        =  local.image_id
  flavor_name     = "s7.large.2"
  use_floatings   = "diable"
  login_mode      = "key_pair"
  key_pair_id     = ctyun_keypair.scaling_test.id
  monitor_service = true
  az_names        = ["cn-huadong1-jsnj1A-public-ctcloud"]
  volumes         = [{"volume_type":"SATA", "volume_size": 40, "flag":"OS"}]
}


resource "ctyun_scaling_group" "group_test" {
  security_group_id_list = [ctyun_security_group.sg_test.id]
  name                   = "scaling-group-for-policy"
  health_mode            = "server"
  subnet_id_list         = [ctyun_subnet.subnet_test.id]
  move_out_strategy      = "earlier_config"
  vpc_id                 = ctyun_vpc.vpc_test.id
  min_count              = 1
  max_count              = 3
  health_period          = 300
  use_lb                 = 2
  config_list            = [ctyun_scaling_config.config_test.id]
  az_strategy            = "uniform_distribution"
  delete_protection      = "disable"
}

# 告警策略
resource "ctyun_scaling_policy" "policy_alert_example" {
  group_id                    = ctyun_security_group.sg_test.id
  name                        = "alert-policy-example"
  policy_type                 = "alert"
  operate_unit                = "percent"
  operate_count               = 10
  action                      = "increase"
  cooldown                    = 300
  trigger_name                = "cpu-high-alert"
  trigger_metric_name         = "cpu_util"
  trigger_statistics          = "avg"
  trigger_comparison_operator = "ge"
  trigger_threshold           = 80
  trigger_period              = "5m"
  trigger_evaluation_count    = 3
}
