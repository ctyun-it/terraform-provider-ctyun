resource "ctyun_acl_rule" "%[1]s" {
  acl_id                 = "%[2]s"
  direction              = "%[3]s"
  protocol               = "%[4]s"
  ip_version             = "%[5]s"
  source_ip_address      = "%[6]s"
  destination_ip_address = "%[7]s"
  action                 = "%[8]s"
  enabled                = %[9]t
  description            = "%[10]s"
  priority               = %[11]d
}
