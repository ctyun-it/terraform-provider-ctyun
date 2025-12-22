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


# 创建vpc
resource "ctyun_vpc" "vpc_test" {
  name        = "vpc-for-ecs"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
}

# 在vpc下创建子网
resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "subnet-for-ecs"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns = [
    "114.114.114.114",
    "8.8.8.8",
  ]
}

# 查询可用镜像
data "ctyun_images" "image_test" {
  name       = "CentOS Linux 8.4"
  visibility = "public"
  page_no    = 1
  page_size  = 10
}

# 查询可用规格
data "ctyun_ecs_flavors" "ecs_flavor_test" {
  cpu    = 2
  ram    = 4
  arch   = "x86"
  series = "S"
  type   = "CPU_S7"
}

# 导入密钥对
resource "ctyun_keypair" "keypair_test" {
  name       = "keypair-for-ecs"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjUnAnTid4wmVtajSmElMtH03OvOyY81ybfswbUu9Gt83DVVzDnwb3rcQW1us8SeKm/gRINkgdrRAgfXAmTKR7AorYtWWc/tzb6kcDpL2E8Qk+n6cyFAxXNoX2vXBr4kC9wz1uwjGyxoSlpHLIpscfI0Ef652gMlSyfODehAJHj3JPMr8pvtPIUqsZI3JOGTUzxaA2JVC0LxQegphYYf2TxGd9GLRUv1p/0BUAPCMg1NaITXNVEj3A11hk1nrFoJMmvIwIUkLmRuQcxuNAdxeLB7GXXVjKpnKIJL4L64dyA9GWa3Gb7gCJyRaBc5UhK4hT57wmukCrldHHtdF1IJr"
}

resource "ctyun_ecs" "ecs_test" {
  instance_name    = "ecs-demo"
  display_name     = "ecs-demo"
  flavor_id        = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].id
  image_id         = data.ctyun_images.image_test.images[0].id
  system_disk_type = "sata"
  system_disk_size = 100
  vpc_id           = ctyun_vpc.vpc_test.id
  subnet_id        = ctyun_subnet.subnet_test.id
  key_pair_name    = ctyun_keypair.keypair_test.name
  cycle_type       = "on_demand"
}

resource "ctyun_ebs" "ebs_test" {
  count      = 5
  name       = "ebs-data-test${count.index + 1}"
  mode       = "vbd"
  type       = "sata"
  size       = 60
  cycle_type = "on_demand"
}

resource "ctyun_ecs_data_volume" "test" {
  instance_id = ctyun_ecs.ecs_test.id
  ebs_ids     = ctyun_ebs.ebs_test[*].id
  # ebs_ids = [for ebs in ctyun_ebs.ebs_test : ebs.id]
}