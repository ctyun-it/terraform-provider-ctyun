resource "ctyun_mysql_param_template" "%[1]s" {
  name         = "%[2]s"
  engine       = "%[3]s"
  description  = "%[4]s"
  template_parameters = %[5]s
}

