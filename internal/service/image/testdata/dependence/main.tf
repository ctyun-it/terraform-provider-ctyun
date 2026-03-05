resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-image"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-image"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns         = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-image"
  description = "terraform测试使用"
  lifecycle {
    prevent_destroy = true
  }
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

# 创建数据盘资源
resource "ctyun_ebs" "data_disk_test" {
  name       = "tf-test-data-disk"
  mode       = "vbd"
  type       = "sata"
  size       = 60
  cycle_type = "on_demand"
}

# 创建ECS实例资源
resource "ctyun_ecs" "ecs_test" {
  instance_name      = "tf-test-ecs"
  display_name       = "tf-test-init-ecs"
  flavor_id           = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].id
  image_id            = data.ctyun_images.image_test.images[0].id
  system_disk_type   = "sata"
  system_disk_size   = 60
  vpc_id             = ctyun_vpc.vpc_test.id
  subnet_id          = ctyun_subnet.subnet_test.id
  security_group_ids = [ctyun_security_group.security_group_test.id]
  cycle_type         = "on_demand"
}

# 创建EBS与ECS的关联关系（显式挂载）
resource "ctyun_ebs_association_ecs" "data_disk_association" {
  instance_id = ctyun_ecs.ecs_test.id
  ebs_id      = ctyun_ebs.data_disk_test.id
}
