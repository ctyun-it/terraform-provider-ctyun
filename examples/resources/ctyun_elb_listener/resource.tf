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

resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-elb"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-elb"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_elb_loadbalancer" "listener_test" {
  subnet_id     = ctyun_subnet.subnet_test.id
  name          = "tf-elb-for-listener"
  sla_name      = "elb.s2.small"
  resource_type = "internal"
  vpc_id        = ctyun_vpc.vpc_test.id
  cycle_type    = "month"
  cycle_count   = 1
}

resource "ctyun_elb_target_group" "test2" {
  name      = "tf-tg-for-target2"
  vpc_id    = ctyun_vpc.vpc_test.id
  algorithm = "wrr"
}

resource "ctyun_elb_listener" "elb_listener_test" {
  loadbalancer_id     = ctyun_elb_loadbalancer.listener_test.id
  name                = "tf_elb_listener"
  protocol            = "TCP"
  protocol_port       = 12345
  default_action_type = "forward"
  target_groups = [{ target_group_id = ctyun_elb_target_group.test2.id }]
  listener_cps        = 1
  establish_timeout   = 100
}

