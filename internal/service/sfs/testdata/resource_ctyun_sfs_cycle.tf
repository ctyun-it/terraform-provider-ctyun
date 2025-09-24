resource "ctyun_sfs" "%[1]s" {
  is_encrypt   = %[2]t
  sfs_type     = "%[3]s"
  sfs_protocol = "%[4]s"
  name         = "%[5]s"
  sfs_size     = %[6]d
  cycle_type   = "%[7]s"
  cycle_count = %[8]d
  vpc_id       = "%[9]s"
  subnet_id    = "%[10]s"
}



