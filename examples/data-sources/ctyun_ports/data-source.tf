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

data "ctyun_ports" "port_data" {
}

output "ctyun_ports_value" {
  value = data.ctyun_ports.port_data
}