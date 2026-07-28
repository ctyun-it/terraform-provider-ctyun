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

data "ctyun_private_nat_snats" "private_nat_snat" {
  nat_gateway_id = "natgw-ltxyq3aa7z"
}

output "ctyun_private_nat_snats_value" {
  value = data.ctyun_private_nat_snats.private_nat_snat
}
