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

resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-ecs"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-ecs"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns = [
    "114.114.114.114",
    "8.8.8.8",
  ]
}

resource "ctyun_ecs" "test" {
  instance_name    = "tf-test-ecs"
  display_name     = "tf-test-init-ecs"
  flavor_id        = "9b4b5e39-db25-f2c8-3914-76881ee77d5c"
  image_id         = "fa3f3784-34f9-4f6b-80a1-dd173d486bd6"
  system_disk_type = "sata"
  system_disk_size = 60
  vpc_id           = ctyun_vpc.vpc_test.id
  subnet_id        = ctyun_subnet.subnet_test.id
  key_pair_name    = "tf-keypair-for-ecs"
  cycle_type       = "on_demand"
}
