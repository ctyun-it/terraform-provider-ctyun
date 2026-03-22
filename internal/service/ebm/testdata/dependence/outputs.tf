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

output "ebs_az" {
  value =  local.first_cloud_boot_true_az
}

output "raid_az" {
  value = local.first_cloud_boot_false_az
}

output "standard_az" {
  value = local.standard_az
}

output "standard_subnet_id" {
  value = ctyun_subnet.subnet_ebm.id
}

output "standard_device_type" {
  value = local.standard_device_type
}

output "standard_image" {
  value = data.ctyun_ebm_device_images.standard_test.images[0].image_uuid
}

output "standard_system_raid" {
  value = local.standard_system_raid_id
}

output "standard_data_raid" {
  value = local.standard_data_raid_id
}

