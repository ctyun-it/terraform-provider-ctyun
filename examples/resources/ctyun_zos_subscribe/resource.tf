
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

resource "ctyun_zos_subscribe" "example" {
  region_id =  "200000003327"
}