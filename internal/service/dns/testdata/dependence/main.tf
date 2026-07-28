resource "ctyun_vpc" "vpc_test" {
  count = 6
  name        = "tf-vpc-for-dns-${count.index+1}"
  cidr        = "192.168.0.0/16"
  description = "terraform-iaas测试使用"
  enable_ipv6 = true
}

resource "ctyun_private_zone" "zone_test" {
name          = "zone.test.com"
description   = "terraform前置资源"
proxy_pattern = "zone"
ttl           = 300
vpc_id_list   = [ctyun_vpc.vpc_test[0].id, ctyun_vpc.vpc_test[1].id]
}







