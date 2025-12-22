resource "ctyun_acl_rule" "%[1]s" {
  project_id             = "%[2]s"
  acl_id                 = "%[3]s"
  direction              = "%[4]s"
  protocol               = "%[5]s"
  ip_version             = "%[6]s"
  source_port            = "%[7]s"
  destination_port       = "%[8]s"
  source_ip_address      = "%[9]s"
  destination_ip_address = "%[10]s"
  action                 = "%[11]s"
  enabled                = %[12]t
  description            = "%[13]s"
  priority               = %[14]d
}
