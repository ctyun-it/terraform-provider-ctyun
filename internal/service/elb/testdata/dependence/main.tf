// main.tf负责创建或查询单测依赖的前置资源
resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-elb"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-elb"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
}


resource "ctyun_elb_loadbalancer" "test" {
  subnet_id     = ctyun_subnet.subnet_test.id
  name          = "tf-elb-for-rule"
  sla_name      = "elb.s2.small"
  resource_type = "internal"
  vpc_id        = ctyun_vpc.vpc_test.id
  cycle_type    = "on_demand"
}
#
resource "ctyun_elb_loadbalancer" "listener_test" {
  subnet_id     = ctyun_subnet.subnet_test.id
  name          = "tf-elb-for-listener"
  sla_name      = "elb.s2.small"
  resource_type = "internal"
  vpc_id        = ctyun_vpc.vpc_test.id
  cycle_type    = "month"
  cycle_count   = 1
}

resource "ctyun_elb_health_check" "test" {
  name     = "tf-hc-for-targetgroup12"
  protocol = "TCP"
}

resource "ctyun_elb_target_group" "test1" {
  name      = "tf-tg-for-target1_12"
  vpc_id    = ctyun_vpc.vpc_test.id
  algorithm = "wrr"
}

resource "ctyun_elb_target_group" "test2" {
  name      = "tf-tg-for-target2_12"
  vpc_id    = ctyun_vpc.vpc_test.id
  algorithm = "wrr"
}

resource "ctyun_elb_target_group" "test3" {
  name      = "tf-tg-for-target3_3"
  vpc_id    = ctyun_vpc.vpc_test.id
  algorithm = "wrr"
}

resource "ctyun_elb_target_group" "test4" {
  name      = "tf-tg-for-target2_12"
  vpc_id    = ctyun_vpc.vpc_test.id
  algorithm = "wrr"
  protocol  = "HTTP"
  session_sticky_mode = "REWRITE"
  rewrite_cookie_name = "cookies"
}

resource "ctyun_elb_listener" "test" {
  loadbalancer_id     = ctyun_elb_loadbalancer.test.id
  name                = "tf-listener-for-rule"
  protocol            = "HTTP"
  protocol_port       = 12345
  default_action_type = "forward"
  target_groups = [{ target_group_id : ctyun_elb_target_group.test1.id }]
}

data "ctyun_images" "image_test" {
  name       = "CentOS Linux 8.4"
  visibility = "public"
  page_no = 1
  page_size = 10
}

data "ctyun_ecs_flavors" "ecs_flavor_test" {
  cpu    = 2
  ram    = 4
  arch   = "x86"
  series = "C"
  type   = "CPU_C7"
}

#
resource "ctyun_ecs" "ecs_test" {
  instance_name       = "tf-ecs-for-elb"
  display_name        = "tf-ecs-for-elb"
  flavor_id           = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].id
  image_id            = data.ctyun_images.image_test.images[0].id
  system_disk_type    = "sata"
  system_disk_size    = 40
  vpc_id = ctyun_vpc.vpc_test.id
  password            = var.password
  cycle_type          = "on_demand"
  subnet_id = ctyun_subnet.subnet_test.id
  is_destroy_instance = false
}

variable "password" {
  type      = string
  sensitive = true
}