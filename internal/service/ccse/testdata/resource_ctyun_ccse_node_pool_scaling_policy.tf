resource "ctyun_ccse_node_pool_scaling_policy" "%[1]s" {
  cluster_id = "%[2]s"
  values_yaml = <<EOF
%[3]sEOF
}
