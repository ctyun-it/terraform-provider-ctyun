resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-service"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-service"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
  enable_ipv6 = true
}

resource "ctyun_ecs" "ecs_test2" {
  instance_name       = "tf-ecs-for-0709"
  display_name        = "tf-ecs-for-0709"
  flavor_id           = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].id
  image_id            = data.ctyun_images.image_test.images[0].id
  system_disk_type    = "sata"
  system_disk_size    = 40
  vpc_id = ctyun_vpc.vpc_test.id
  password            = var.password
  cycle_type          = "on_demand"
  subnet_id = ctyun_subnet.subnet_test.id
}

resource "ctyun_vpc" "vpc_test2" {
  name        = "tf-vpc-vpce"
  cidr        = "192.167.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test2" {
  vpc_id = ctyun_vpc.vpc_test2.id
  name        = "tf-subnet-vpce2"
  cidr        = "192.167.1.0/24"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
  enable_ipv6 = true
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
  instance_name       = "tf-ecs-for-0627-1"
  display_name        = "tf-ecs-for-0627-1"
  flavor_id           = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].id
  image_id            = data.ctyun_images.image_test.images[0].id
  system_disk_type    = "sata"
  system_disk_size    = 40
  vpc_id = ctyun_vpc.vpc_test2.id
  password            = var.password
  cycle_type          = "on_demand"
  subnet_id = ctyun_subnet.subnet_test2.id
}


resource "ctyun_vpce_service" "test" {
  name  = "tf-vpce-server-0709"
  vpc_id = ctyun_vpc.vpc_test.id
  subnet_id = ctyun_subnet.subnet_test.id
  auto_connection = true
  type = "interface"
  instance_type = "vm"
  instance_id = ctyun_ecs.ecs_test2.id
  rules = [{
    protocol = "TCP"
    endpoint_port = 123
    server_port = 456
  },
  ]
}

resource "ctyun_vpce" "test" {
  name  = "tf-vpce-0709"
  endpoint_service_id = ctyun_vpce_service.test.id
  vpc_id = ctyun_vpc.vpc_test2.id
  subnet_id = ctyun_subnet.subnet_test2.id
  whitelist_flag = false
}

variable "password" {
  type      = string
  sensitive = true
}