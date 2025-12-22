resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-vpc"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-vpc"
  cidr        = "192.168.2.0/24"
  description = "terraform测试使用"
  dns = [
    "114.114.114.114",
    "8.8.8.8",
  ]
  enable_ipv6 = true
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-vpc"
  description = "terraform测试使用"
}

resource "ctyun_bandwidth" "bandwidth_test" {
  name       = "tf-bandwidth-for-vpc"
  bandwidth  = 5
  cycle_type = "on_demand"
}

resource "ctyun_eip" "eip_test" {
  name                = "tf-eip-for-vpc"
  bandwidth           = 5
  cycle_type          = "on_demand"
  demand_billing_type = "upflowc"
}

data "ctyun_images" "image_test" {
  name       = "CentOS Linux 8.4"
  visibility = "public"
  page_no    = 1
  page_size  = 10
}

data "ctyun_ecs_flavors" "ecs_flavor_test" {
  cpu    = 2
  ram    = 4
  arch   = "x86"
  series = "C"
  type   = "CPU_C7"
}

resource "ctyun_ecs" "ecs_test" {
  instance_name    = "tf-ecs-for-vpc"
  display_name     = "tf-ecs-for-vpc"
  flavor_id        = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].id
  image_id         = data.ctyun_images.image_test.images[0].id
  system_disk_type = "sata"
  system_disk_size = 40
  vpc_id           = ctyun_vpc.vpc_test.id
  password         = var.password
  cycle_type       = "on_demand"
  subnet_id        = ctyun_subnet.subnet_test.id
}

variable "password" {
  type      = string
  sensitive = true
}

resource "ctyun_port" "port_test" {
  name                       = "tf-test-port"
  description                = "tf-test-port"
  subnet_id                  = ctyun_subnet.subnet_test.id
  secondary_private_ip_count = 1
}

resource "ctyun_ecs_port_association" "port_association" {
  instance_id = ctyun_ecs.ecs_test.id
  port_id     = ctyun_port.port_test.id
}

resource "ctyun_vip" "vip_test" {
  subnet_id  = ctyun_subnet.subnet_test.id
  vpc_id     = ctyun_vpc.vpc_test.id
  ip_address = "192.168.2.152"
  vip_type   = "v4"
}

resource "ctyun_vpc" "vpc_test2" {
  name        = "tf-vpc-for-vpc2"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_dhcpoptionset" "test" {
  name         = "tf_dhcpoption_set"
  description  = "test"
  domain_name  = "www.ctyun.cn"
  dns_list     = ["114.114.114.114"]
}