output "vpc_id" {
  value = ctyun_vpc.vpc_test.id
}

output "subnet_id" {
  value = ctyun_subnet.subnet_test.id
}

output "security_group_id" {
  value = ctyun_security_group.security_group_test.id
}

output "mongodb_id"{
  value = ctyun_mongodb_instance.mongodb_eip.id
}

output "mongodb_host_ip"{
  value = ctyun_mongodb_instance.mongodb_eip.host_ip
}

output "eip_id" {
  value = ctyun_eip.eip_test.id
}

output "az_name" {
  value = data.ctyun_zones.az.zones[0]
}

output "flavor_name" {
  value = local.flavor_name
}

output "flavor_name2" {
  value = local.flavor_name2
}