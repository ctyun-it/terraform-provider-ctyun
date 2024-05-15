terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
}

data "ctyun_regions" "ctyun_regions_test" {
  name = "南昌"
}

output "ctyun_regions_test" {
  value = data.ctyun_regions.ctyun_regions_test.regions
}
