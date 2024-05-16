terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

resource "ctyun_eip_association" "eip_association_test1" {
  eip_id      = "eip-p9qvl63yt6"
  instance_id = "0b9897c8-ff01-42b4-c4c2-6a427d8b2e9a"
}

resource "ctyun_eip_association" "eip_association_test2" {
  eip_id      = "eip-nl78g1t31o"
  instance_id = "fd94fbe2-26b2-5dbb-5deb-65b4167ca28e"
  region_id   = "200000002527"
  project_id  = "4f5ef15300724760af59b37cf6409f45"
}