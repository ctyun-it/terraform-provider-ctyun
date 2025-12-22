resource "ctyun_mysql_rds_parameter_template" "%[1]s" {
  instance_id     = "%[2]s"
  project_id  = "%[3]s"
  template_id = %[4]d
}