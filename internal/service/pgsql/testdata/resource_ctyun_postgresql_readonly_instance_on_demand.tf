resource "ctyun_postgresql_readonly_instance" "%[1]s" {
  instance_id     = "%[2]s"
  cycle_type  = "%[3]s"
  flavor_name = "%[4]s"
  name        = "%[5]s"
}
