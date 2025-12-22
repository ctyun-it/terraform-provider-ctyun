
output "vpc_id" {
  value = local.data_vpc_id == "" ? ctyun_vpc.vpc_test[0].id : local.data_vpc_id
}

output "subnet_id" {
  value = local.data_subnet_id == "" ? ctyun_subnet.subnet_test[0].id : local.data_subnet_id
}

output "security_group_id" {
  value = local.data_security_group_id == "" ? ctyun_security_group.security_group_test[0].id : local.data_security_group_id
}

output "instance_id" {
  value = local.ecs_instance_id == "" ? ctyun_ecs.ecs_test[0].id : local.ecs_instance_id
}

output "data_disk_id" {
  value = local.data_disk_id == "" ? ctyun_ebs.data_disk_test[0].id : local.data_disk_id
}
