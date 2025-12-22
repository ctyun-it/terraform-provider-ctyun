# 云主机
data "ctyun_ccse_images" "%[1]s" {
  instance_type = "ecs"
  flavor_name = "%[2]s"
}


# 物理机
data "ctyun_ccse_images" "%[3]s" {
  instance_type = "ebm"
  flavor_name = "%[4]s"
}

