resource "ctyun_ec_cloud_gateway" "%[1]s" {
  ec_id       = "%[2]s"
  name        = "%[3]s"
  description = "%[4]s"
  region_id   = "200000002401"
  region_name = "cn-hn-cs42-hncs1A-public-ctcloud"
}
