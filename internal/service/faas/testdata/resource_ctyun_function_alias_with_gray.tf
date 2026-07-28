resource "ctyun_function_alias" "%[1]s" {
  function_name   = "%[2]s"
  alias_name      = "alias-gray-%[1]s"
  version_id      = "%[3]s"
  description     = "%[5]s"
  gray_version_id = "%[4]s"
  gray_weight     = 10
  depends_on = [%[6]s]
}
