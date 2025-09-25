resource "ctyun_ccse_node_association" "%[1]s" {
  cluster_id = "%[2]s"
  instance_type = "%[3]s"
  instance_id = "%[4]s"
  mirror_id = "%[5]s"
  visibility_post_host_script = "%[6]s"
  visibility_host_script = "%[7]s"
  password = "%[8]s"
}