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

resource "ctyun_prefix_list" "example" {
  name         = "prefix-example"
  description  = "prefix list example"
  limit        = 100
  address_type = "ipv6"
  prefix_list_rules = [
    { "cidr" : "2001:db8::/32", "description" : "IPv6 rule 1" },
    { "cidr" : "2001:db8:1::/48", "description" : "IPv6 rule 2" },
  ]
}
