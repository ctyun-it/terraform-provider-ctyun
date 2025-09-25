terraform {
	required_providers {
		ctyun = {
			source = "ctyun-it/ctyun"
		}
	}
}

provider "ctyun" {
	env = "prod"
}

resource "ctyun_vpc" "vpc_test" {
	name        = "tf-vpc-for-nat"
	cidr        = "192.168.0.0/16"
	description = "terraform测试使用"
	enable_ipv6 = true
}

resource "ctyun_nat" "nat_test"{
	vpc_id = ctyun_vpc.vpc_test.id
	spec = 1
	name = "tf-nat"
	description = "terraform测试使用"
	cycle_type = "on_demand"
}

resource "ctyun_nat" "nat_cycle_test" {
	vpc_id = ctyun_vpc.vpc_test.id
	spec = 1
	name = "tf-nat-cycle"
	description = "terraform测试使用"
	cycle_type = "month"
	cycle_count = 1
}
