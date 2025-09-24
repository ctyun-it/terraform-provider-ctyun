output "vpc_id" {
  value = ctyun_vpc.vpc_test.id
}

output "subnet_id" {
  value = ctyun_subnet.subnet_test.id
}

output "security_group_id" {
  value = ctyun_security_group.security_group_test.id
}

output "security_group_id2" {
  value = ctyun_security_group.security_group_test2.id
}

output "device_type" {
  value = local.device_type1
}

output "system_raid" {
  value  = local.system_raid_id
}

output "data_raid" {
  value  = local.data_raid_id
}

output "image_uuid" {
  value = data.ctyun_ebm_device_images.test.images[0].image_uuid
}

output "ebs_id" {
  value = ctyun_ebs.ebs_test.id
}

output "ebm_id" {
  value = ctyun_ebm.ebm_test.id
}

output "az2" {
  value = local.az2
}