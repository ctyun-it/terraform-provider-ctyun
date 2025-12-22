resource "ctyun_postgresql_readonly_instance" "%[1]s" {
  instance_id     = "%[2]s"
  cycle_type  = "%[3]s"
  flavor_name = "%[4]s"
  project_id  = "%[5]s"
  name        = "%[6]s"
}
