resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-vpce"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-vpce"
  cidr        = "192.168.1.0/24"
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
  instance_name       = "tf-ecs-for-vpce"
  display_name        = "tf-ecs-for-vpce"
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

resource "ctyun_ecs" "ecs_test2" {
  instance_name       = "tf-ecs-for-vpce2"
  display_name        = "tf-ecs-for-vpce2"
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


resource "ctyun_vpce_service" "vpce_service_test" {
  name  = "tf-vpce-server-for-vpce"
  vpc_id = ctyun_vpc.vpc_test.id
  subnet_id = ctyun_subnet.subnet_test.id
  auto_connection = true
  type = "interface"
  instance_id = ctyun_ecs.ecs_test.id
  instance_type = "vm"
  rules = [{
    protocol = "TCP"
    endpoint_port = 999
    server_port = 999
  }]
}

resource "ctyun_vpce_service" "reverse_vpce_service_test" {
  name  = "tf-reverse-vpce-server-for-vpce"
  vpc_id = ctyun_vpc.vpc_test.id
  subnet_id = ctyun_subnet.subnet_test.id
  auto_connection = true
  type = "reverse"
}

data "ctyun_vpce_service_transit_ips" "vpce_service_transit_ip_test" {
  endpoint_service_id = ctyun_vpce_service.reverse_vpce_service_test.id
}

resource "ctyun_vpce" "vpce_test" {
  name  = "tf-vpce-for-vpce"
  endpoint_service_id = ctyun_vpce_service.reverse_vpce_service_test.id
  vpc_id = ctyun_vpc.vpc_test.id
  subnet_id = ctyun_subnet.subnet_test.id
  whitelist_flag = false
}

variable "password" {
  type      = string
  sensitive = true
}