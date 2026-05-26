
resource "ctyun_function_alias" "%[1]s" {
  function_name = "%[2]s"
  alias_name    = "%[3]s"
  version_id    = "%[4]s"
  description   = "%[5]s"
  depends_on = [%[6]s]
}
