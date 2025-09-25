resource "ctyun_elb_acl" "%[1]s" {
  name = "%[2]s"
  source_ips = [%[3]s]
}
