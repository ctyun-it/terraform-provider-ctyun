resource "ctyun_oceanfs" "%[1]s" {
  protocol = "%[2]s"
  name         = "%[3]s"
  size     = "%[4]d"
  cycle_type   = "%[5]s"
  cycle_count  = "%[6]d"
  vpc_id       = "%[7]s"
  subnet_id    = "%[8]s"
  tags         = [%[9]s]
}

