resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-cloud-assistant"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id     = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-cloud-assistant"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8"
  ]
}

data "ctyun_images" "image_test" {
  name       = "CentOS Linux 8.4"
  visibility = "public"
  page_no    = 1
  page_size  = 10
}

data "ctyun_ecs_flavors" "ecs_flavor_test" {
  cpu  = 2
  ram  = 4
  arch = "x86"
}

resource "ctyun_ecs" "ecs_test1" {
  instance_name       = "tf-ecs-for-cloud-assistant1"
  display_name        = "tf-ecs-for-cloud-assistant1"
  flavor_id           = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].id
  image_id            = data.ctyun_images.image_test.images[0].id
  system_disk_type    = "SATA"
  system_disk_size    = 40
  vpc_id              = ctyun_vpc.vpc_test.id
  password            = var.password
  cycle_type          = "on_demand"
  subnet_id           = ctyun_subnet.subnet_test.id
  is_destroy_instance = false
}

resource "ctyun_ecs" "ecs_test2" {
  instance_name       = "tf-ecs-for-cloud-assistant2"
  display_name        = "tf-ecs-for-cloud-assistant2"
  flavor_id           = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].id
  image_id            = data.ctyun_images.image_test.images[0].id
  system_disk_type    = "SATA"
  system_disk_size    = 40
  vpc_id              = ctyun_vpc.vpc_test.id
  password            = var.password
  cycle_type          = "on_demand"
  subnet_id           = ctyun_subnet.subnet_test.id
  is_destroy_instance = false
}


# 管理云助手命令
resource "ctyun_cloud_assistant_command" "command_test" {
  command_name    = "tf-command-test"
  command_type    = "shell"
  command_content = "echo hello world"
  description     = "云助手前置资源"
  timeout         = 1000
}



variable "password" {
  type      = string
  sensitive = true
}
