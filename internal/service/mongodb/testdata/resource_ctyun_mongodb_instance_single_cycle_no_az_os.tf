resource "ctyun_mongodb_instance" "%[1]s" {
  cycle_type             = "%[2]s"
  cycle_count            = %[3]d
  vpc_id                 = "%[4]s"
  flavor_name            = "%[5]s"
  subnet_id              = "%[6]s"
  security_group_id      = "%[7]s"
  name                   = "%[8]s"
  password               = "%[9]s"
  prod_id                = "%[10]s"
  read_port              = %[11]d
  storage_type           = "%[12]s"
  storage_space          = %[13]d
  backup_storage_type    = "%[14]s"  // 对象存储
}
