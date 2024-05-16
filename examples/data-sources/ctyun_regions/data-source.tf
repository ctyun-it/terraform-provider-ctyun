terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

data "ctyun_regions" "ctyun_regions_test" {
  name = "南昌"
}

output "ctyun_regions_test" {
  value = data.ctyun_regions.ctyun_regions_test.regions
}
