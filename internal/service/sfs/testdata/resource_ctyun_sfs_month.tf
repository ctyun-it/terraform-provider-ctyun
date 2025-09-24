resource "ctyun_sfs" "%[1]s" {
  is_encrypt   = %[2]t
  kms_uuid     = "%[3]s"
  sfs_type     = "%[4]s"
  sfs_protocol = "%[5]s"
  name         = "%[6]s"
  sfs_size     = %[7]d
  cycle_type   = "%[8]s"
  cycle_count = %[9]d
  vpc_id       = "%[10]s"
  subnet_id    = "%[11]s"
  read_only    = %[12]t
}


