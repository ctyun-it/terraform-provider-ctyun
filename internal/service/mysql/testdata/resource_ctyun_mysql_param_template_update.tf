resource "ctyun_mysql_param_template" "%[1]s" {
  project_id   = "%[2]s"
  name         = "%[3]s"
  engine       = "%[4]s"
  description  = "%[5]s"
  template_parameters = %[6]s
}