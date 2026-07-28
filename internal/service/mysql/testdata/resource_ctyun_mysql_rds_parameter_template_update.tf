resource "ctyun_mysql_rds_parameter_template" "%[1]s" {
  instance_id    = "%[2]s"
  parameters = %[3]s
}