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

output "template_id" {
  value = tostring(data.ctyun_mysql_param_templates.template.param_templates[0].id)
}

output "task_id" {
  value = data.ctyun_mysql_backups.backup_test.backups.0.records.0.task_id
}

output "backup_timestamp" {
  value = data.ctyun_mysql_recoverable_time_points.time_point_test.backup_time_points.0.end_time
}
