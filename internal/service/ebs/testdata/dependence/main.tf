resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-ebs"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-ebs"
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
  instance_name       = "tf-ecs-for-ebs"
  display_name        = "tf-ecs-for-ebs"
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

resource "ctyun_ebs" "ebs_test" {
  name       = "ebs-test25"
  mode       = "vbd"
  type       = "sata"
  size       = 60
  cycle_type = "on_demand"
}

resource "ctyun_ebs_snapshot" "ebs_snapshot_test" {
  name = "ebs_snapshot_test"
  disk_id = ctyun_ebs.ebs_test.id
  retention_policy = "forever"
}

variable "password" {
  type      = string
  sensitive = true
}
