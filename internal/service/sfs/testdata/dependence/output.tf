output "vpc_id" {
  value = ctyun_vpc.vpc_test.id
}

output "vpc_id1" {
  value = ctyun_vpc.vpc_test1.id
}

output "subnet_id" {
  value = ctyun_subnet.subnet_test.id
}

output "sfs_uid" {
  value = ctyun_sfs.sfs_test.id
}

output "sfs_permission_group_id" {
  value = ctyun_sfs_permission_group.group_test.id
}

output "sfs_permission_group_id1" {
  value = ctyun_sfs_permission_group.group_test1.id
}