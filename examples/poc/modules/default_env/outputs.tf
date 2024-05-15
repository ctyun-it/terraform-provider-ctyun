output "default_vpc_id" {
  value = module.vpc.vpc_id
}

output "default_subnet_id" {
  value = module.subnet.subnet_id
}

output "default_security_group_id" {
  value = module.security_group.security_group_id
}

output "default_key_pair_name" {
  value = module.key_pair.key_pair_name
}