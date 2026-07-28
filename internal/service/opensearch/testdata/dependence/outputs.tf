output "region_id" {
  value = ctyun_vpc.vpc_test.region_id
}

output "vpc_id" {
  value = ctyun_vpc.vpc_test.id
}

output "subnet_id" {
  value = ctyun_subnet.subnet_test.id
}

output "security_group_id" {
  value = ctyun_security_group.security_group_test.id
}

output "az_name" {
  value = local.az_name
}

output "zone_list" {
  value = data.ctyun_zones.test.zones[0]
}

output "flavor_name" {
  value = local.flavor_name
}

output "storage_type" {
  value = local.storage_type
}

