resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-gaokeyong"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-vpc-gaokeyong"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
  enable_ipv6 = true
}

resource "ctyun_eip" "eip_test" {
  name                = "tf-eip-test134"
  bandwidth           = 1
  cycle_type          = "on_demand"
  demand_billing_type = "bandwidth"
}

data "ctyun_zones" "test" {

}

locals {
  az1 = data.ctyun_zones.test.zones[0]
  az2 = data.ctyun_zones.test.zones[1]
}

output "az1" {
  value = local.az1
}

data "ctyun_images" "image_test" {
  az_name    = local.az1
  name       = "CentOS Linux 8.4"
  visibility = "public"
  page_no    = 1
  page_size  = 10
}

data "ctyun_ecs_flavors" "ecs_flavor_test" {
  az_name = local.az1
  cpu     = 2
  ram     = 4
  arch    = "x86"
  series  = "C"
  type    = "CPU_C7"
}

resource "ctyun_ecs" "ecs_test" {
  instance_name    = "tf-ecs-gaokeyong-1"
  display_name     = "tf-ecs-gaokeyong-1"
  flavor_id        = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].id
  image_id         = data.ctyun_images.image_test.images[0].id
  system_disk_type = "sata"
  system_disk_size = 40
  vpc_id           = ctyun_vpc.vpc_test.id
  password         = var.password
  az_name          = local.az1
  cycle_type       = "on_demand"
  subnet_id        = ctyun_subnet.subnet_test.id
}

# az2
data "ctyun_images" "image_test2" {
  az_name    = local.az2
  name       = "CentOS Linux 8.4"
  visibility = "public"
  page_no    = 1
  page_size  = 10
}

data "ctyun_ecs_flavors" "ecs_flavor_test2" {
  az_name = local.az2
  cpu     = 2
  ram     = 4
  arch    = "x86"
  series  = "C"
  type    = "CPU_C7"
}

resource "ctyun_ecs" "ecs_test2" {
  instance_name    = "tf-ecs-gaokeyong-2"
  display_name     = "tf-ecs-gaokeyong-2"
  flavor_id        = data.ctyun_ecs_flavors.ecs_flavor_test2.flavors[0].id
  image_id         = data.ctyun_images.image_test2.images[0].id
  system_disk_type = "sata"
  system_disk_size = 40
  vpc_id           = ctyun_vpc.vpc_test.id
  password         = var.password
  az_name          = local.az2
  cycle_type       = "on_demand"
  subnet_id        = ctyun_subnet.subnet_test.id
}

resource "ctyun_elb_loadbalancer" "test" {
  az_name       = local.az1
  subnet_id     = ctyun_subnet.subnet_test.id
  name          = "tf-elb-gaokeyong"
  sla_name      = "elb.s2.small"
  resource_type = "external"
  vpc_id        = ctyun_vpc.vpc_test.id
  eip_id        = ctyun_eip.eip_test.id
  cycle_type    = "month"
  cycle_count   = 1
}

resource "ctyun_elb_listener" "elb_listener_test" {
  loadbalancer_id     = ctyun_elb_loadbalancer.test.id
  name                = "tf-elb-listener-gaokeyong"
  protocol            = "TCP"
  protocol_port       = 456
  default_action_type = "forward"
  target_groups = [{ target_group_id = ctyun_elb_target_group.target_group_test.id }]
}

resource "ctyun_elb_target_group" "target_group_test" {
  name      = "tf-target-group"
  vpc_id    = ctyun_vpc.vpc_test.id
  algorithm = "wrr"
}

resource "ctyun_elb_target" "target" {
  target_group_id = ctyun_elb_target_group.target_group_test.id
  instance_type = "VM"
  instance_id = ctyun_ecs.ecs_test.id
  protocol_port = 456
}

resource "ctyun_elb_target" "target2" {
  target_group_id = ctyun_elb_target_group.target_group_test.id
  instance_type = "VM"
  instance_id = ctyun_ecs.ecs_test2.id
  protocol_port = 456
}

variable "password" {
  type      = string
  sensitive = true
}