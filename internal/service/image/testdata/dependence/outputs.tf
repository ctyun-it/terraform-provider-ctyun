
output "vpc_id" {
  value = ctyun_vpc.vpc_test.id
}

output "subnet_id" {
  value = ctyun_subnet.subnet_test.id
}

output "security_group_id" {
  value = ctyun_security_group.security_group_test.id
}

output "instance_id" {
  value = ctyun_ecs.ecs_test.id
}

output "data_disk_id" {
  value = ctyun_ebs.data_disk_test.id
}
