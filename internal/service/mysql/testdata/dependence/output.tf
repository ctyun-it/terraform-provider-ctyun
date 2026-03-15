output "vpc_id" {
  value = ctyun_vpc.vpc_test.id
}

output "subnet_id" {
  value = ctyun_subnet.subnet_test.id
}

output "security_group_id" {
  value = ctyun_security_group.security_group_test.id
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

output "template_id" {
  value = tostring(data.ctyun_mysql_param_templates.template.param_templates[0].id)
}

output "task_id" {
  value = data.ctyun_mysql_backups.backup_test.backups.0.records.0.task_id
}

output "backup_timestamp" {
  value = data.ctyun_mysql_recoverable_time_points.time_point_test.backup_time_points.0.end_time
}

output "flavor_name" {
  value = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].name
}

output "flavor_name2" {
  value = data.ctyun_ecs_flavors.ecs_flavor_test2.flavors[0].name
}