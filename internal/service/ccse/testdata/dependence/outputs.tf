output "vpc_id" {
  value = local.real_vpc_id
}

output "subnet_id" {
  value = local.real_subnet_id
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

output "ebm_id" {
  value = ctyun_ebm.ebm_test.id
}

output "ecs_mirror_id" {
  value = "3d2c356a-685a-4e8c-b904-bb0725bfc220"
}

output "ebm_mirror_id" {
  value =  "im-lplf1yqhl3mewvc5pjvha70wklej"
}

output "device_type" {
  value =  local.device_type1
}

output "ebm_mirror_name" {
  value = "CTyunOS23.01@cpu_ccse_img_4.0_09"
}

output "ebm_az" {
  value = local.az2
}