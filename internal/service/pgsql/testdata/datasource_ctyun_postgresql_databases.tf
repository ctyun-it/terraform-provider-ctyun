data "ctyun_postgresql_databases" "%[1]s" {
  instance_id = "%[2]s"
  name        = "%[3]s"
}
