terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

resource "ctyun_bandwidth_association_eip" "bandwidth_association_eip_test" {
  bandwidth_id = "bandwidth-at2yy664m5"
  eip_id       = "eip-p9qvl63yt6"
}