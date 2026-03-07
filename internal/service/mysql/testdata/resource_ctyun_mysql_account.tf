resource "ctyun_mysql_account" "%[1]s" {
  instance_id          = "%[2]s"
  name     = "%[3]s"
  password = "%[4]s"
  schema_privilege_list = %[5]s
  description      = "%[6]s"
}


