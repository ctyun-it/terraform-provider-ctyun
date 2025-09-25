output "vpc_id" {
  value = ctyun_vpc.vpc_test.id
}

output "subnet_id" {
  value = ctyun_subnet.subnet_test.id
}
# output "loadbalancer_id" {
#   value = ctyun_elb_loadbalancer.test.id
# }
output "loadbalancer_id_rule" {
  value = ctyun_elb_loadbalancer.listener_test.id
}

output "health_check_id" {
  value = ctyun_elb_health_check.test.id
}

output "target_group_id" {
  value = ctyun_elb_target_group.test1.id
}
output "target_group_id2" {
  value = ctyun_elb_target_group.test2.id
}
output "target_group_id3" {
  value = ctyun_elb_target_group.test3.id
}
output "target_group_id4" {
  value = ctyun_elb_target_group.test4.id
}
output "listener_id" {
  value = ctyun_elb_listener.test.id
}

output "instance_id" {
  value = ctyun_ecs.ecs_test.id
}