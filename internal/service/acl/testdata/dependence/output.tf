output "vpc_id" {
  value = ctyun_vpc.vpc_test.id
}

output "subnet_id" {
  value = ctyun_subnet.subnet_test[0].id
}
output "subnet_id2" {
  value = ctyun_subnet.subnet_test[1].id
}

output "acl_id" {
  value = ctyun_acl.acl_test.id
}

output "acl_id2" {
  value = ctyun_acl.acl_subnet_test.id
}