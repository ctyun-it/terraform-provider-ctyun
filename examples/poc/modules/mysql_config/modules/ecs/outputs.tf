output "ecs_instance_id" {
  value = ctyun_ecs.mysql_test[*].id
}