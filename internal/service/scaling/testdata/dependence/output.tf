
output "vpc_id" {
  value = local.real_vpc_id
}

output "subnet_id" {
  value = local.real_subnet_id
}

output "subnet_id1"{
  value = ctyun_subnet.subnet_test1.id
}

output "security_group_id" {
  value = local.real_security_group_id
}
output "security_group_id1" {
  value = ctyun_security_group.security_group_test1.id
}

output "key_pair_id" {
  value = ctyun_keypair.scaling_test.id
}
output "elb_loadbalancer_id" {
  value = ctyun_elb_loadbalancer.elb_test.id
}
output "elb_target_group_id" {
  value = ctyun_elb_target_group.target_group_test.id
}

output "elb_loadbalancer_id1" {
  value = ctyun_elb_loadbalancer.elb_test1.id
}
output "elb_target_group_id1" {
  value = ctyun_elb_target_group.target_group_test1.id
}

output "scaling_group_id" {
  value = tostring(ctyun_scaling_group.group_test.id)
}

output "scaling_config_id" {
  value = tostring(ctyun_scaling_config.config_test.id)
}

output "scaling_config_id1" {
  value = tostring(ctyun_scaling_config.config_test1.id)
}

output "instance_uuid" {
  value = ctyun_ecs.ecs_test.id
}

output "instance_uuid1" {
  value = ctyun_ecs.ecs_test1.id
}

output "instance_uuid2" {
  value = ctyun_ecs.ecs_test2.id
}

output "instance_uuid3" {
  value = ctyun_ecs.ecs_test3.id
}