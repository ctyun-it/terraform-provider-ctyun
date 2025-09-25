output "vpc_id" {
  value = local.real_vpc_id
}

output "subnet_id" {
  value = local.real_subnet_id
}

output "security_group_id" {
  value = local.real_security_group_id
}

output "eip_id" {
  value = ctyun_eip.eip_test.id
}

output "mysql_id" {
  value = ctyun_mysql_instance.mysql_test.id
}

output "az_name" {
  value = local.az_name
}