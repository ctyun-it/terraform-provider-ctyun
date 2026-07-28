output "vpc_id" {
  value = local.real_vpc_id
}

output "subnet_id" {
  value =  local.real_subnet_id
}

output "security_group_id" {
  value = local.real_security_group_id
}

output "rocketmq_single_spec_name" {
  value = local.single_spec_name
}

output "rocketmq_single_spec_name2" {
  value = local.single_spec_name2
}

output "rocketmq_cluster_spec_name" {
  value = local.cluster_spec_name
}

output "rocketmq_cluster_spec_name2" {
  value = local.cluster_spec_name2
}
#
output "zone" {
  value = data.ctyun_zones.test.zones[0]
}

output "instance_id" {
  value = ctyun_rocketmq_instance.example.id
}
