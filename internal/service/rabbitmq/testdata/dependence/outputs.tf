output "vpc_id" {
  value = local.real_vpc_id
}

output "subnet_id" {
  value =  local.real_subnet_id
}

output "security_group_id" {
  value = local.real_security_group_id
}

output "rabbitmq_single_disk_type" {
  value = local.single_disk_type
}

output "rabbitmq_single_spec_name" {
  value = local.single_spec_name
}

output "rabbitmq_single_spec_name2" {
  value = local.single_spec_name2
}

output "rabbitmq_cluster_disk_type" {
  value = local.cluster_disk_type
}

output "rabbitmq_cluster_spec_name" {
  value = local.cluster_spec_name
}

output "rabbitmq_cluster_spec_name2" {
  value = local.cluster_spec_name2
}