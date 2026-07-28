resource "ctyun_sfs" "%[1]s" {
  is_encrypt   = %[2]t
  type     = "%[3]s"
  protocol = "%[4]s"
  name         = "%[5]s"
  size     = %[6]d
  cycle_type   = "%[7]s"
  cycle_count = %[8]d
  vpc_id       = "%[9]s"
  subnet_id    = "%[10]s"
}



