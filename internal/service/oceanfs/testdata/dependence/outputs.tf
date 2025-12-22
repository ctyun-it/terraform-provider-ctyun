output "vpc_id" {
  value = ctyun_vpc.vpc_test.id
}

output "vpc_id1" {
  value = ctyun_vpc.vpc_test1.id
}

output "subnet_id" {
  value = ctyun_subnet.subnet_test.id
}

output "subnet_id1" {
  value = ctyun_subnet.subnet_test1.id
}

output "permission_group_id" {
  value = ctyun_oceanfs_permission_group.test.id
}

output "permission_group_id1" {
  value = ctyun_oceanfs_permission_group.test1.id
}

output "oceanfs_id" {
  value = ctyun_oceanfs.test.id
}