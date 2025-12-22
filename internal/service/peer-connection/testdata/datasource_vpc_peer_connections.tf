data "ctyun_vpc_peer_connections" "%[1]s" {
  provider = ctyun.region
}
provider "ctyun" {
  alias     = "region"
  # region_id = "200000002689"
  env       = "prod"
  ak        = "dd2befb766fe423bb3d360563ee786d3"
  sk        = "7cf01ea846d64a7e9c48acb9d06b120e"
}