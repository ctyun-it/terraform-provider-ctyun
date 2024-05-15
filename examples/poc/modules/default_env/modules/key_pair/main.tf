terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
}

resource "ctyun_keypair" "keypair_test" {
  name       = var.key_pair_name
  public_key = var.public_key
}