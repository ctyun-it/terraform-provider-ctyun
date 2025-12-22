resource "ctyun_postgresql_account" "%[1]s" {
  project_id = "%[2]s"
  instance_id = "%[3]s"
  name = "%[4]s"
  password = "%[5]s"
  user_type = "%[6]s"
  schema_privilege_list = %[7]s
  description = "%[8]s"
}
