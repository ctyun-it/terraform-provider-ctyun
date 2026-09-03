resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-ai"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-ai"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns = [
    "114.114.114.114",
    "8.8.8.8",
  ]
  enable_ipv6 = true
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-ai"
  description = "terraform测试使用"
}

data "ctyun_ecs_flavors" "ecs_flavor_test" {
  cpu  = 2
  ram  = 4
  arch = "x86"
}

locals {
  available_flavor = [for f in data.ctyun_ecs_flavors.ecs_flavor_test.flavors : f if f.available == true][0]
}

data "ctyun_images" "image_test" {
  name        = "CentOS Linux 8.4"
  visibility  = "public"
  page_no     = 1
  page_size   = 10
  flavor_name = local.available_flavor.name
}

resource "ctyun_ecs" "ecs_test" {
  instance_name    = "tf-ecs-for-ai"
  display_name     = "tf-ecs-for-ai"
  flavor_id        = local.available_flavor.id
  image_id         = data.ctyun_images.image_test.images[0].id
  system_disk_type = "SATA"
  system_disk_size = 40
  vpc_id           = ctyun_vpc.vpc_test.id
  password         = var.password
  cycle_count      = 2
  auto_renew       = false
  cycle_type       = "on_demand"
  subnet_id        = ctyun_subnet.subnet_test.id
}

resource "ctyun_ecs" "ecs_test2" {
  instance_name    = "tf-ecs-for-ai2"
  display_name     = "tf-ecs-for-ai2"
  flavor_id        = local.available_flavor.id
  image_id         = data.ctyun_images.image_test.images[0].id
  system_disk_type = "SATA"
  system_disk_size = 40
  vpc_id           = ctyun_vpc.vpc_test.id
  password         = var.password
  cycle_count      = 1
  auto_renew       = true
  cycle_type       = "on_demand"
  subnet_id        = ctyun_subnet.subnet_test.id
}

resource "ctyun_ebs" "ebs_test" {
  name        = "ai-data-volume"
  mode        = "vbd"
  type        = "SATA"
  size        = 60
  cycle_type  = "on_demand"
  cycle_count = 2
}

variable "password" {
  type      = string
  sensitive = true
}

resource "ctyun_eip" "eip_test" {
  name                = "tf-eip-for-ai"
  bandwidth           = 5
  cycle_type          = "on_demand"
  cycle_count         = 2
  demand_billing_type = "upflowc"
}

resource "ctyun_eip" "eip_test2" {
  name                = "tf-eip-for-ai2"
  bandwidth           = 5
  cycle_type          = "month"
  cycle_count         = 1
  demand_billing_type = "bandwidth"
}

resource "ctyun_bandwidth" "bandwidth_test" {
  name        = "tf-bandwidth-for-ai"
  bandwidth   = 5
  cycle_type  = "on_demand"
  cycle_count = 2
}

resource "ctyun_ipv6_bandwidth" "ipv6_bandwidth_test" {
  name        = "ipv6-bandwidth-for-ai"
  bandwidth   = "5"
  cycle_type  = "on_demand"
  cycle_count = 2
}

resource "ctyun_nat" "nat_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  spec        = 1
  name        = "tf-nat-for-ai"
  cycle_type  = "on_demand"
  cycle_count = 2
}

resource "ctyun_private_nat" "nat_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  spec        = "small"
  subnet_id   = ctyun_subnet.subnet_test.id
  name        = "tf-private-for-ai"
  cycle_type  = "on_demand"
  cycle_count = 2
}

resource "ctyun_elb_loadbalancer" "test" {
  subnet_id     = ctyun_subnet.subnet_test.id
  name          = "tf-elb-for-ai"
  sla_name      = "elb.s2.small"
  resource_type = "internal"
  vpc_id        = ctyun_vpc.vpc_test.id
  cycle_type    = "on_demand"
  cycle_count   = 2
}
