resource "ctyun_hpfs" "%[1]s" {
  protocol = "%[2]s"
  name     = "%[3]s"
  size     = %[4]d
  cluster_name = "%[5]s"
  baseline     = "%[6]s"
}

