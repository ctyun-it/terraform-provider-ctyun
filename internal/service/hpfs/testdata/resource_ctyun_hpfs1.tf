resource "ctyun_hpfs" "%[1]s" {
  sfs_protocol = "%[2]s"
  name     = "%[3]s"
  sfs_size     = %[4]d
  az_name      = "%[5]s"
  cluster_name = "%[6]s"
  baseline     = "%[7]s"
}

