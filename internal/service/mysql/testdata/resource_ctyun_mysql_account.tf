resource "ctyun_mysql_account" "%[1]s" {
  instance_id          = "%[2]s"
  project_id       = "%[3]s"
  name     = "%[4]s"
  password = "%[5]s"
  schema_privilege_list = %[6]s
  description      = "%[7]s"
}


