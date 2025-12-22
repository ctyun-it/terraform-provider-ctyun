resource "ctyun_acl_rule" "%[1]s" {
  project_id             = "%[2]s"
  acl_id                 = "%[3]s"
  direction              = "%[4]s"
  protocol               = "%[5]s"
  ip_version             = "%[6]s"
  source_ip_address      = "%[7]s"
  destination_ip_address = "%[8]s"
  action                 = "%[9]s"
  enabled                = %[10]t
  description            = "%[11]s"
  priority               = %[12]d
}
