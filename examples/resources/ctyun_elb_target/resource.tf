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
resource "ctyun_elb_target_group" "target_group_test" {
  name      = "tf-tg-for-target1_12"
  vpc_id    = ctyun_vpc.vpc_test.id
  algorithm = "wrr"
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

resource "ctyun_elb_target" "elb_target_test" {
  target_group_id = ctyun_elb_target_group.target_group_test.id
  instance_type = "VM"
  instance_id = ctyun_ecs.ecs_test.id
  protocol_port = 12345
}

