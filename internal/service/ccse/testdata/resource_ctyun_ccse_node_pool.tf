resource "ctyun_ccse_node_pool" "%[1]s" {
  name           = "%[2]s"
  cycle_type              = "%[9]s"
  %[10]s
  visibility_post_host_script = "%[3]s"
  visibility_host_script = "%[4]s"
  instance_type            = "ecs"
  mirror_id                = "3f80d8c0-8eb5-4afa-a506-13ba68b61872"
  mirror_type              = 1
  password                 = "%[14]s"
  max_pod_num              = 110
  item_def_name            = "%[11]s"
  cluster_id               = "%[12]s"
  node_num  = %[13]d
  sys_disk = {
    type = "%[5]s"
    size = %[6]d
  }

  data_disks = [
    {
      type = "%[7]s"
      size = %[8]d
    }
  ]
  az_infos = [
    {
      az_name = "cn-huadong1-jsnj1A-public-ctcloud"
    }
  ]
}
