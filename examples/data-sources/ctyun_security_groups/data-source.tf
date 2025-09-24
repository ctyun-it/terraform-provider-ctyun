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
  name        = "vpc-test-mc22"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
  region_id   = "200000001852"
}

resource "ctyun_security_group" "security_group_test" {
  region_id = "200000001852"
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "terraform-minchiang"
  description = "terraform测试使用"
}


data "ctyun_security_groups" "test" {
  region_id = "200000001852"
  security_group_id = ctyun_security_group.security_group_test.id
  # page_no = 1
  # page_size = 1
}

output "ctyun_test" {
  value = data.ctyun_security_groups.test
}

