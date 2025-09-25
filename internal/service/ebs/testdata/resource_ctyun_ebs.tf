resource "ctyun_ebs" "%[1]s" {
  name       = "%[2]s"
  mode       = "vbd"
  type       = "sata"
  size       = %[3]d
  cycle_type = "on_demand"
}
