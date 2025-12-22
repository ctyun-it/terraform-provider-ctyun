resource "ctyun_ebs" "%[1]s" {
  name       = "%[2]s"
  mode       = "vbd"
  type       = "%[3]s"
  size       = %[4]d
  cycle_type = "on_demand"
  %[5]s
}
