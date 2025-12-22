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

data "ctyun_private_nat_transit_ips" "transit_ip" {
  nat_gateway_id = "natgw-ltxyq3aa7z"
}

output "ctyun_private_nat_transit_ips_value" {
  value = data.ctyun_private_nat_transit_ips.transit_ip
}