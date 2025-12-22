data "ctyun_zones" "test" {

}

locals {
  az_name    = data.ctyun_zones.test.zones[0]

}
data "ctyun_hpfs_clusters" "test" {

}
