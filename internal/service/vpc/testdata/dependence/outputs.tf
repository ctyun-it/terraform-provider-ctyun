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

output "network_interface_id" {
  value = ctyun_port.port_test.id
}
output "subnet_id" {
  value = ctyun_subnet.subnet_test.id
}

output "vip_id" {
  value = ctyun_vip.vip_test.id
}

output "vpc_id2" {
  value = ctyun_vpc.vpc_test2.id
}

output "dhcp_id" {
  value = ctyun_dhcpoptionset.test.id
}