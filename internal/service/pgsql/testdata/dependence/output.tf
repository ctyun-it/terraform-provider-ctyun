output "vpc_id" {
  value = ctyun_vpc.vpc_test.id
}

output "subnet_id" {
  value = ctyun_subnet.subnet_test.id
}

output "security_group_id1" {
  value = ctyun_security_group.security_group_test1.id
}
output "security_group_id2" {
  value = ctyun_security_group.security_group_test2.id
}

output "security_group_id3" {
  value = ctyun_security_group.security_group_test3.id
}

output "eip_id" {
  value = ctyun_eip.eip_test.id
}

output "pgsql_id" {
  value = ctyun_postgresql_instance.test.id
}

output "az_name" {
  value = data.ctyun_zones.az.zones[0]
}

output "param_template_id"{
  value =  tostring(data.ctyun_postgresql_param_templates.param_templates.parameter_templates[0].id)
}

output "charset_name" {
  value = data.ctyun_postgresql_character_set.charsets.character_set[1]
}

output "collate_name" {
  value = data.ctyun_postgresql_collation_time_zone.collations.collations[0].coll_name
}
output "collate_type" {
  value = data.ctyun_postgresql_collation_time_zone.collations.collations[0].coll_type
}

output "account_name" {
  value = data.ctyun_postgresql_accounts.accounts.accounts[0].name
}


output "flavor_name" {
  value = local.flavor_name
}

output "flavor_name2" {
  value = local.flavor_name2
}