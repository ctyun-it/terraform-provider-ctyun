resource "ctyun_sdwan_acl" "acl_test" {
  name = "acl_test"

  rules = [{
    direction        = "in"
    protocol         = "udp"
    ip_version       = "IPv4"
    dst_cidr         = "10.0.0.0/16"
    dst_port_range   = "-1/-1"
    priority         = 100
    action           = "allow"
    src_cidr         = "10.0.0.0/16"
    src_port_range   = "-1/-1"
  }]
}

