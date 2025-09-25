resource "ctyun_sfs" "%[1]s" {
  sfs_type     = "%[2]s"
  sfs_protocol = "%[3]s"
  name         = "%[4]s"
  sfs_size     = %[5]d
  cycle_type   = "%[6]s"
  az_name      = "%[7]s"
  vpc_id       = "%[8]s"
  subnet_id    = "%[9]s"
  read_only    = %[10]t
}

