output "vpc_id" {
  value = ctyun_vpc.vpc_test.id
}

output "subnet_id" {
  value = ctyun_subnet.subnet_test.id
}

output "ecs_id" {
  value = ctyun_ecs.ecs_test.id
}

output "ecs_id2" {
  value = ctyun_ecs.ecs_test2.id
}

output "vpce_service_id" {
  value = ctyun_vpce_service.vpce_service_test.id
}

output "reverse_vpce_service_id" {
  value = ctyun_vpce_service.reverse_vpce_service_test.id
}

output "vpce_id" {
  value = ctyun_vpce.vpce_test.id
}

output "ecs_fixed_ip" {
  value = ctyun_ecs.ecs_test.fixed_ip
}

output "vpce_service_transit_ip" {
  value = data.ctyun_vpce_service_transit_ips.vpce_service_transit_ip_test.ips[0].transit_ip
}