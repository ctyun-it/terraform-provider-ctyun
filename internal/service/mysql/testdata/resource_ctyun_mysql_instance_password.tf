resource "ctyun_mysql_instance" "%[1]s" {
  cycle_type        = "%[2]s"
  vpc_id            = "%[3]s"
  subnet_id         = "%[4]s"
  security_group_id = "%[5]s"
  name              = "%[6]s"
  password          = "%[7]s"
  flavor_name       = "%[8]s"
  prod_id           = "%[9]s"
  storage_type      = "%[10]s"
  storage_space     = %[11]d
}

