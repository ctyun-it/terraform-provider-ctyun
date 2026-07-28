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

resource "ctyun_dhcpoptionset" "example" {
  name        = "dhcpoptionset"
  description = "Example DHCP option set"
  domain_name = "example.com"
  dns_list    = ["8.8.8.8", "114.114.114.114"]
}