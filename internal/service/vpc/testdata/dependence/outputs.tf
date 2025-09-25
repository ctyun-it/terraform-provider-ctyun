output "vpc_id" {
  value = ctyun_vpc.vpc_test.id
}

output "bandwidth_id" {
  value = ctyun_bandwidth.bandwidth_test.id
}

output "eip_id" {
  value = ctyun_eip.eip_test.id
}

output "ecs_id" {
  value = ctyun_ecs.ecs_test.id
}

output "security_group_id" {
  value = ctyun_security_group.security_group_test.id
}