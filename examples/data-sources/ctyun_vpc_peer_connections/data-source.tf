terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

data "ctyun_vpc_peer_connections" "test" {

}
