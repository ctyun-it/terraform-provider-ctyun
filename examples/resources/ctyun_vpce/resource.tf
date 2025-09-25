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
  name        = "tf-vpc-test-qqq"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-test"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
  enable_ipv6 = true
}

resource "ctyun_vpce_service" "test" {
  name  = "tf-vpce-server-sss"
  vpc_id = ctyun_vpc.vpc_test.id
  subnet_id = ctyun_subnet.subnet_test.id
  auto_connection = true
  type = "interface"
  instance_id = "d40b78e2-23de-4fa6-baf0-e500750f985b"
  instance_type = "vm"
  rules = [{
    protocol = "TCP"
    endpoint_port = 2
    server_port = 2
  },
  ]
}

resource "ctyun_vpce" "test" {
  name  = "tf-vpce-123"
  endpoint_service_id = ctyun_vpce_service.test.id
  vpc_id = ctyun_vpc.vpc_test.id
  subnet_id = ctyun_subnet.subnet_test.id
  whitelist_flag = true
  whitelist_cidr = ["192.168.1.0/24"]
}