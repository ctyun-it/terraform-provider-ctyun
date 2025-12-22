resource "ctyun_postgresql_instance" "%[1]s" {
  cycle_type            = "%[2]s"
  prod_id               = "%[3]s"
  flavor_name           = "%[4]s"
  storage_type          = "%[5]s"
  storage_space         = %[6]d
  name                  = "%[7]s"
  password              = "%[8]s"
  case_sensitive        = %[9]t
  vpc_id                = "%[10]s"
  subnet_id             = "%[11]s"
  security_group_id     = "%[12]s"
  backup_storage_type  = "%[13]s"
}