resource "ctyun_postgresql_param_template" "%[1]s" {
  name = "%[2]s"
  source_template_id = %[3]d
  description = "%[4]s"
  template_parameters = %[5]s
}
