data "ctyun_services" "test" {
  type = "region"
}

data "ctyun_regions" "test" {
  name = "华北2"
}

data "ctyun_zones" "test" {
  region_id = data.ctyun_regions.test.regions[0].id
}