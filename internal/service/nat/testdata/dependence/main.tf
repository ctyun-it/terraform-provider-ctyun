resource "ctyun_eip" "eip_test" {
  name                = "tf-eip-for-nat"
  bandwidth           = 1
  cycle_type          = "on_demand"
  demand_billing_type = "upflowc"
}

resource "ctyun_eip" "eip_test1" {
  name                = "tf-eip-for-nat1"
  bandwidth           = 1
  cycle_type          = "on_demand"
  demand_billing_type = "upflowc"
}

resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-nat"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_nat" "nat_test"{
  vpc_id = ctyun_vpc.vpc_test.id
  spec = 1
  name = "tf-nat-for-test"
  description = "terraform测试使用"
  cycle_type = "on_demand"
}

resource "ctyun_subnet" "subnet_test1" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-nat1"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_subnet" "subnet_test2" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-nat2"
  cidr        = "192.168.128.0/24"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
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

resource "ctyun_ecs" "ecs_test" {
  instance_name       = "tf-ecs-for-nat"
  display_name        = "tf-ecs-for-nat"
  flavor_id           = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].id
  image_id            = data.ctyun_images.image_test.images[0].id
  system_disk_type    = "sata"
  system_disk_size    = 40
  vpc_id = ctyun_vpc.vpc_test.id
  password            = var.password
  cycle_type          = "on_demand"
  subnet_id = ctyun_subnet.subnet_test1.id
  is_destroy_instance = false
}

variable "password" {
  type      = string
  sensitive = true
}