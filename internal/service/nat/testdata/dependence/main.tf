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
  name        = "tf-vpc-for-private-nat"
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
    "8.8.8.8"
  ]
}

resource "ctyun_subnet" "subnet_test2" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-nat2"
  cidr        = "192.168.128.0/24"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8"
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
resource "ctyun_private_nat" "nat_test"{
  vpc_id = ctyun_vpc.vpc_test.id
  spec = "small"
  subnet_id = ctyun_subnet.subnet_test2.id
  name = "tf-private_nat"
  description = "terraform测试使用"
  cycle_type = "on_demand"
}

resource "ctyun_private_nat_transit_ip" "ip1"{
  nat_gateway_id = ctyun_private_nat.nat_test.id
  address ="192.168.128.100"
}

resource "ctyun_private_nat_transit_ip" "ip2"{
  nat_gateway_id = ctyun_private_nat.nat_test.id
  address ="192.168.128.101"
}


variable "password" {
  type      = string
  sensitive = true
}
resource "ctyun_security_group" "security_group_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-private-nat"
  description = "terraform测试使用"
}

resource "ctyun_port" "port" {
  name                       = "port-test-update"
  description                = "port 测试-测试"
  subnet_id                  = ctyun_subnet.subnet_test1.id
  security_group_ids        =  [ctyun_security_group.security_group_test.id]
  secondary_private_ip_count = 1
}

resource "ctyun_ecs_port_association" "ecs_port_for_association_test" {
  instance_id          =  ctyun_ecs.ecs_test.id
  port_id = ctyun_port.port.id
}