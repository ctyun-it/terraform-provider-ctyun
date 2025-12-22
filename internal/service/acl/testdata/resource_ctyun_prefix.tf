resource "ctyun_prefix_list" "%[1]s" {
  name              = "%[2]s"
  description       = "%[3]s"
  limit             = %[4]d
  address_type      = "%[5]s"
  prefix_list_rules = [%[6]s]
}
