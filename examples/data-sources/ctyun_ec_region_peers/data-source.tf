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


data "ctyun_ec_region_peers" "peers_examples" {
  ec_id = "49410d6d-fd53-48b3-9f78-cb28da38d7be"
}
