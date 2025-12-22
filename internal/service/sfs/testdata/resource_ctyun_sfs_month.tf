resource "ctyun_sfs" "%[1]s" {
  is_encrypt   = %[2]t
  kms_uuid     = "%[3]s"
  type     = "%[4]s"
  protocol = "%[5]s"
  name         = "%[6]s"
  size     = %[7]d
  cycle_type   = "%[8]s"
  cycle_count = %[9]d
  vpc_id       = "%[10]s"
  subnet_id    = "%[11]s"
}


