resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-acl"
  cidr        = "192.168.0.0/16"
  description = "terraform-iaas测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  count       = 2
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-acl-${count.index+1}"
  cidr        = "192.168.${count.index+1}.0/24"
  description = "terraform测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_acl" "acl_test" {
  project_id  = "0"
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "acl-terraform-test"
  description = "terraform测试使用"
}

resource "ctyun_acl" "acl_subnet_test" {
  project_id  = "0"
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "acl-subnet-terraform-test"
  description = "terraform测试使用"
}
