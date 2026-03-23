output "vpc_id" {
  value = local.real_vpc_id
}

output "subnet_id" {
  value = local.real_subnet_id
}


output "subnet_id2" {
  value = local.real_subnet_id2
}

output "security_group_id" {
  value = ctyun_security_group.security_group_test.id
}

output "flavor_name" {
  value = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].name
}

output "cluster_id" {
  value = ctyun_ccse_cluster.test.id
}

output "chart_name" {
  value = local.chart_name
}

output "chart_version1" {
  value = local.chart_version1
}

output "chart_version2" {
  value = local.chart_version2
}

output "chart_values_yaml" {
  value = jsonencode(data.ctyun_ccse_plugin_market.test1.values)
}

output "chart_values_json" {
  value = jsonencode(data.ctyun_ccse_plugin_market.test2.values)
}

output "ecs_id" {
  value = ctyun_ecs.ecs_test.id
}

output "ecs_mirror_id" {
  value = lookup([for img in data.ctyun_ccse_images.ccse_images.images : img if img.os_distro == "CTyunOS"][0], "id")
}

# output "ebm_id" {
#   value = ctyun_ebm.ebm_test.id
# }
#
# output "ebm_mirror_id" {
#   value =  "im-lplf1yqhl3mewvc5pjvha70wklej"
# }
#
output "device_type" {
  value =  local.first_cloud_boot_false.device_type
}

output "ebm_az" {
  value =  local.first_cloud_boot_false.az_name
}

#
# output "ebm_mirror_name" {
#   value = "CTyunOS23.01@cpu_ccse_img_4.0_09"
# }
#

output "ecs_az" {
  value = ctyun_ecs.ecs_test.az_name
}