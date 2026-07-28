resource "ctyun_oceanfs" "%[1]s" {
  project_id   = "%[2]s"
  protocol = "%[3]s"
  name         = "%[4]s"
  size     = %[5]d
  cycle_type   = "%[6]s"
  vpc_id       = "%[7]s"
  subnet_id    = "%[8]s"
}

