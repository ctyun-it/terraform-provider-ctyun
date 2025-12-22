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

data "ctyun_private_nats" "private_nat_data" {
}

output "ctyun_private_nats_value" {
  value = data.ctyun_private_nats.private_nat_data
}