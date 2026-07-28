terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 可参考index.md，在环境变量中配置ak、sk、资源池ID、可用区名称
provider "ctyun" {
  env = "prod"
}

resource "ctyun_dhcpoptionset" "example" {
  name        = "dhcpoptionset"
  description = "Example DHCP option set"
  domain_name = "example.com"
  dns_list    = ["8.8.8.8", "114.114.114.114"]
}

data "ctyun_dhcpoptionset_association_vpcs" "example" {
  dhcp_option_sets_id = ctyun_dhcpoptionset.example.id
}

output "ctyun_dhcpoptionset_association_vpcs_example" {
  value = data.ctyun_dhcpoptionset_association_vpcs.example
}