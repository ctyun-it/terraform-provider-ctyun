resource "ctyun_postgresql_param_template" "%[1]s" {
  project_id = "%[2]s"
  name = "%[3]s"
  source_template_id = %[4]d
  description = "%[5]s"
}
