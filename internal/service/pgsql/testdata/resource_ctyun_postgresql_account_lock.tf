resource "ctyun_postgresql_account" "%[1]s" {
  instance_id = "%[2]s"
  name = "%[3]s"
  password = "%[4]s"
  user_type = "%[5]s"
  schema_privilege_list = %[6]s
  description = "%[7]s"
  is_lock = %[8]t
}
