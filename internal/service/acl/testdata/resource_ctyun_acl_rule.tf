resource "ctyun_acl_rule" "%[1]s" {
  acl_id                 = "%[2]s"
  direction              = "%[3]s"
  protocol               = "%[4]s"
  ip_version             = "%[5]s"
  source_port            = "%[6]s"
  destination_port       = "%[7]s"
  source_ip_address      = "%[8]s"
  destination_ip_address = "%[9]s"
  action                 = "%[10]s"
  enabled                = %[11]t
  description            = "%[12]s"
}
