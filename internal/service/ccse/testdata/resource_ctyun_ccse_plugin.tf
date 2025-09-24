resource "ctyun_ccse_plugin" "%[1]s" {
  cluster_id = "%[2]s"
  chart_name = "%[3]s"
  chart_version = "%[4]s"
  %[5]s
}
