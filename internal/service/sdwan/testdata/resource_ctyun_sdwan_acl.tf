resource "ctyun_sdwan_acl" "%[1]s" {
  name = "%[2]s"

  rules = [{
    direction        = "%[3]s"
    protocol         = "%[4]s"
    ip_version       = "%[5]s"
    dst_cidr         = "%[6]s"
    dst_port_range   = "%[7]s"
    priority         = %[8]d
    action           = "%[9]s"
    src_cidr         = "%[10]s"
    src_port_range   = "%[11]s"
  }]
}