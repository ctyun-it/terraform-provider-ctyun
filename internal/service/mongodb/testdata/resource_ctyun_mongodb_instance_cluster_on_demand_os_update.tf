resource "ctyun_mongodb_instance" "%[1]s" {
  cycle_type             = "%[2]s"
  vpc_id                 = "%[3]s"
  flavor_name            = "%[4]s"
  subnet_id              = "%[5]s"
  security_group_id      = "%[6]s"
  name                   = "%[7]s"
  password               = "%[8]s"
  prod_id                = "%[9]s"
  read_port              = %[10]d
  storage_type           = "%[11]s"
  storage_space          = %[12]d
  backup_storage_type    = "%[13]s"
  availability_zone_info = %[14]s
  shard_num              = %[15]d
  mongos_num             = %[16]d
  upgrade_node_type      = "%[17]s"
}
